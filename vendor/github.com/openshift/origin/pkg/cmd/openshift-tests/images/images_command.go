package images

import (
	"fmt"
	"os"
	"sort"
	"strings"

	k8simage "k8s.io/kubernetes/test/utils/image"

	"github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/origin/pkg/clioptions/imagesetup"
	"github.com/openshift/origin/pkg/cmd"
	"github.com/openshift/origin/test/extended/util/image"
	"github.com/spf13/cobra"
	"k8s.io/kube-openapi/pkg/util/sets"
	"k8s.io/kubectl/pkg/util/templates"
)

func NewImagesCommand() *cobra.Command {
	o := &imagesOptions{}
	cmd := &cobra.Command{
		Use:   "images",
		Short: "Gather images required for testing",
		Long: templates.LongDesc(fmt.Sprintf(`
		Creates a mapping to mirror test images to a private registry

		This command identifies the locations of all test images referenced by the test
		suite and outputs a mirror list for use with 'oc image mirror' to copy those images
		to a private registry. The list may be passed via file or standard input.

				$ openshift-tests images --to-repository private.com/test/repository > /tmp/mirror
				$ oc image mirror -f /tmp/mirror

		The 'run' and 'run-upgrade' subcommands accept '--from-repository' which will source
		required test images from your mirror.

		See the help for 'oc image mirror' for more about mirroring to disk or consult the docs
		for mirroring offline. You may use a file:// prefix in your '--to-repository', but when
		mirroring from disk to your offline repository you will have to construct the appropriate
		disk to internal registry statements yourself.

		By default, the test images are sourced from a public container image repository at
		%[1]s and are provided as-is for testing purposes only. Images are mirrored by the project
		to the public repository periodically.
		`, imagesetup.DefaultTestImageMirrorLocation)),
		PersistentPreRun: cmd.NoPrintVersion,
		SilenceUsage:     true,
		SilenceErrors:    true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := imagesetup.VerifyTestImageRepoEnvVarUnset(); err != nil {
				return err
			}

			if o.Verify {
				return imagesetup.VerifyImages()
			}

			repository := o.Repository
			var prefix string
			for _, validPrefix := range []string{"file://", "s3://"} {
				if strings.HasPrefix(repository, validPrefix) {
					repository = strings.TrimPrefix(repository, validPrefix)
					prefix = validPrefix
					break
				}
			}
			ref, err := reference.Parse(repository)
			if err != nil {
				return fmt.Errorf("--to-repository is not valid: %v", err)
			}
			if len(ref.Tag) > 0 || len(ref.ID) > 0 {
				return fmt.Errorf("--to-repository may not include a tag or image digest")
			}

			if err := imagesetup.VerifyImages(); err != nil {
				return err
			}
			lines, err := createImageMirrorForInternalImages(prefix, ref, !o.Upstream)
			if err != nil {
				return err
			}
			for _, line := range lines {
				fmt.Fprintln(os.Stdout, line)
			}
			// TODO: these should be removed when landing k8s 1.31:
			newImages := injectNewImages(ref, o.Upstream)
			for _, line := range newImages {
				fmt.Fprintln(os.Stdout, line)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&o.Upstream, "upstream", o.Upstream, "Retrieve images from the default upstream location")
	cmd.Flags().StringVar(&o.Repository, "to-repository", o.Repository, "A container image repository to mirror to.")
	// this is a private flag for debugging only
	cmd.Flags().BoolVar(&o.Verify, "verify", o.Verify, "Verify the contents of the image mappings")
	cmd.Flags().MarkHidden("verify")
	return cmd
}

func injectNewImages(ref reference.DockerImageReference, upstream bool) []string {
	lines := []string{}
	for original, mirror := range map[string]string{
		"registry.k8s.io/etcd:3.5.15-0":                                   "e2e-11-registry-k8s-io-etcd-3-5-15-0-W7c5qq4cz4EE20EQ",
		"registry.k8s.io/e2e-test-images/agnhost:2.52":                    "e2e-1-registry-k8s-io-e2e-test-images-agnhost-2-52-vo_U710PrYLetnfE",
		"registry.k8s.io/pause:3.10":                                      "e2e-27-registry-k8s-io-pause-3-10-b3MYAwZ_MelO9baY",
		"registry.k8s.io/e2e-test-images/regression-issue-74839:1.2":      "e2e-30-registry-k8s-io-e2e-test-images-regression-issue-74839-1-2-pZ_RxNuqvcwEiCKE",
		"registry.k8s.io/e2e-test-images/resource-consumer:1.13":          "e2e-31-registry-k8s-io-e2e-test-images-resource-consumer-1-13-LT0C2W4wMzShSeGS",
		"registry.k8s.io/e2e-test-images/volume/nfs:1.4":                  "e2e-32-registry-k8s-io-e2e-test-images-volume-nfs-1-4-u7V8iW5QIcWM2i6h",
		"registry.k8s.io/sig-storage/hostpathplugin:v1.14.0":              "e2e-43-registry-k8s-io-sig-storage-hostpathplugin-v1-14-0-LWjla55lyZB4CQu0",
		"registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.10.1":   "e2e-45-registry-k8s-io-sig-storage-csi-node-driver-registrar-v2-10-1-bVz-v06gRSvh6Rp3",
		"registry.k8s.io/sig-storage/csi-attacher:v4.6.1":                 "e2e-47-registry-k8s-io-sig-storage-csi-attacher-v4-6-1-NP4z4EcSo-N1xk_4",
		"registry.k8s.io/sig-storage/csi-provisioner:v5.0.1":              "e2e-48-registry-k8s-io-sig-storage-csi-provisioner-v5-0-1-wPw2vjyYX1LWVmkn",
		"registry.k8s.io/sig-storage/csi-resizer:v1.11.1":                 "e2e-49-registry-k8s-io-sig-storage-csi-resizer-v1-11-1-6jB55ZThgstz1GrW",
		"registry.k8s.io/sig-storage/csi-snapshotter:v8.0.1":              "e2e-50-registry-k8s-io-sig-storage-csi-snapshotter-v8-0-1-vAVT_GHf7Vm-TXyx",
		"registry.k8s.io/e2e-test-images/busybox:1.29-2":                  "e2e-51-registry-k8s-io-e2e-test-images-busybox-1-29-2-ZYWRth-o9U_JR2ZE",
		"registry.k8s.io/sig-storage/hello-populator:v1.0.1":              "e2e-34-registry-k8s-io-sig-storage-hello-populator-v1-0-1-Ei7libli17J5IWn-",
		"registry.k8s.io/sig-storage/volume-data-source-validator:v1.0.0": "e2e-35-registry-k8s-io-sig-storage-volume-data-source-validator-v1-0-0-pJwTeQGTDmAV8753",
		"registry.k8s.io/sig-storage/livenessprobe:v2.12.0":               "e2e-46-registry-k8s-io-sig-storage-livenessprobe-v2-12-0-wCYz5fsB0ew8MCS0",
	} {
		if upstream {
			lines = append(lines, fmt.Sprintf("%s %s:%s", original, ref.Exact(), mirror))
		} else {
			lines = append(lines, fmt.Sprintf("quay.io/openshift/community-e2e-images:%s %s:%s", mirror, ref.Exact(), mirror))
		}
	}
	sort.Strings(lines)
	return lines
}

type imagesOptions struct {
	Repository string
	Upstream   bool
	Verify     bool
}

// createImageMirrorForInternalImages returns a list of 'oc image mirror' mappings from source to
// target or returns an error. If mirrored is true the images are assumed to have already been copied
// from their upstream location into our official mirror, in the REPO:TAG format where TAG is a hash
// of the original internal name and the index of the image in the array. Otherwise the mappings will
// be set to mirror the location as defined in the test code into our official mirror, where the target
// TAG is the hash described above.
func createImageMirrorForInternalImages(prefix string, ref reference.DockerImageReference, mirrored bool) ([]string, error) {
	source := ref.Exact()

	initialDefaults := k8simage.GetOriginalImageConfigs()
	exceptions := image.Exceptions.List()
	defaults := map[k8simage.ImageID]k8simage.Config{}

imageLoop:
	for i, config := range initialDefaults {
		for _, exception := range exceptions {
			if strings.Contains(config.GetE2EImage(), exception) {
				continue imageLoop
			}
		}
		defaults[i] = config
	}

	updated := k8simage.GetMappedImageConfigs(defaults, ref.Exact())
	openshiftDefaults := image.OriginalImages()
	openshiftUpdated := image.GetMappedImages(openshiftDefaults, imagesetup.DefaultTestImageMirrorLocation)

	// if we've mirrored, then the source is going to be our repo, not upstream's
	if mirrored {
		baseRef, err := reference.Parse(imagesetup.DefaultTestImageMirrorLocation)
		if err != nil {
			return nil, fmt.Errorf("invalid default mirror location: %v", err)
		}

		// calculate the mapping of upstream images by setting defaults to baseRef
		covered := sets.NewString()
		for i, config := range updated {
			defaultConfig := defaults[i]
			pullSpec := config.GetE2EImage()
			if pullSpec == defaultConfig.GetE2EImage() {
				continue
			}
			if covered.Has(pullSpec) {
				continue
			}
			covered.Insert(pullSpec)
			e2eRef, err := reference.Parse(pullSpec)
			if err != nil {
				return nil, fmt.Errorf("invalid test image: %s: %v", pullSpec, err)
			}
			if len(e2eRef.Tag) == 0 {
				return nil, fmt.Errorf("invalid test image: %s: no tag", pullSpec)
			}
			config.SetRegistry(baseRef.Registry)
			config.SetName(baseRef.RepositoryName())
			config.SetVersion(e2eRef.Tag)
			defaults[i] = config
		}

		// calculate the mapping for openshift images by populating openshiftUpdated
		openshiftUpdated = make(map[string]string)
		sourceMappings := image.GetMappedImages(openshiftDefaults, imagesetup.DefaultTestImageMirrorLocation)
		targetMappings := image.GetMappedImages(openshiftDefaults, source)

		for from, to := range targetMappings {
			if from == to {
				continue
			}
			if covered.Has(to) {
				continue
			}
			covered.Insert(to)
			from := sourceMappings[from]
			openshiftUpdated[from] = to
		}
	}

	covered := sets.NewString()
	var lines []string
	for i := range updated {
		a, b := defaults[i], updated[i]
		from, to := a.GetE2EImage(), b.GetE2EImage()
		if from == to {
			continue
		}
		if covered.Has(from) {
			continue
		}
		covered.Insert(from)
		lines = append(lines, fmt.Sprintf("%s %s%s", from, prefix, to))
	}

	for from, to := range openshiftUpdated {
		if from == to {
			continue
		}
		if covered.Has(from) {
			continue
		}
		covered.Insert(from)
		lines = append(lines, fmt.Sprintf("%s %s%s", from, prefix, to))
	}

	sort.Strings(lines)
	return lines, nil
}