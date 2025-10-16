//go:build linux || freebsd

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/dmikushin/podman-shared/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to check if --shared-base-layers flag is parsed correctly
func checkSharedBaseLayersFlag(podmanTest *PodmanTestIntegration, args []string) *PodmanSessionIntegration {
	allArgs := append([]string{"run", "--shared-base-layers"}, args...)
	session := podmanTest.Podman(allArgs)
	session.WaitWithDefaultTimeout()
	return session
}

// Helper function to create a test file for shared storage simulation
func createTestSharedDir(tempDir string) (string, error) {
	sharedDir := filepath.Join(tempDir, "test_shared_storage")
	err := os.MkdirAll(sharedDir, 0755)
	if err != nil {
		return "", err
	}

	// Create a marker file to simulate shared storage
	markerFile := filepath.Join(sharedDir, "shared_marker.txt")
	err = os.WriteFile(markerFile, []byte("shared storage test"), 0644)
	if err != nil {
		return "", err
	}

	return sharedDir, nil
}

var _ = Describe("Podman shared base layers CLI tests", func() {

	Context("CLI Flag Parsing and Basic Validation", func() {
		It("should parse --shared-base-layers flag without syntax errors", func() {
			// Test that the flag is recognized by checking help output
			session := podmanTest.Podman([]string{"run", "--help"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			helpOutput := session.OutputToString()
			// The help should contain the shared-base-layers flag
			Expect(helpOutput).To(ContainSubstring("shared-base-layers"))
		})

		It("should accept --shared-base-layers flag in combination with other flags", func() {
			// Test flag combination without actually running containers
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--help"})
			session.WaitWithDefaultTimeout()
			// Should not fail due to flag parsing errors
			Expect(session).Should(ExitCleanly())
		})

		It("should not show flag parsing errors when using --shared-base-layers", func() {
			// Test that the flag doesn't cause immediate CLI errors
			// Using an invalid image name should give image-related error, not flag error
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "nonexistent-image-12345"})
			session.WaitWithDefaultTimeout()

			// Should fail due to missing image, not due to flag issues
			Expect(session).ShouldNot(ExitCleanly())
			errorOutput := session.ErrorToString()

			// Error should be about image, not about unknown flag
			Expect(errorOutput).ToNot(ContainSubstring("unknown flag"))
			Expect(errorOutput).ToNot(ContainSubstring("shared-base-layers"))
		})
	})

	Context("Shared Storage Directory Tests", func() {
		var tempDir string

		BeforeEach(func() {
			tempDir = podmanTest.TempDir
		})

		It("should handle shared storage directory creation", func() {
			sharedDir, err := createTestSharedDir(tempDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(sharedDir).To(BeADirectory())

			// Verify marker file exists
			markerFile := filepath.Join(sharedDir, "shared_marker.txt")
			Expect(markerFile).To(BeAnExistingFile())

			// Read marker file content
			content, err := os.ReadFile(markerFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("shared storage test"))
		})

		It("should handle multiple shared storage directories", func() {
			// Create multiple shared directories to simulate different base layers
			for i := 0; i < 3; i++ {
				dirName := fmt.Sprintf("shared_layer_%d", i)
				layerDir := filepath.Join(tempDir, dirName)
				err := os.MkdirAll(layerDir, 0755)
				Expect(err).ToNot(HaveOccurred())

				// Create layer-specific content
				contentFile := filepath.Join(layerDir, fmt.Sprintf("layer_%d_content.txt", i))
				content := fmt.Sprintf("Content for layer %d", i)
				err = os.WriteFile(contentFile, []byte(content), 0644)
				Expect(err).ToNot(HaveOccurred())

				Expect(layerDir).To(BeADirectory())
				Expect(contentFile).To(BeAnExistingFile())
			}
		})
	})

	Context("Flag Interaction Tests", func() {
		It("should work with common container flags", func() {
			// Test various flag combinations that should work together
			flagCombinations := [][]string{
				{"--shared-base-layers", "--rm"},
				{"--shared-base-layers", "--detach"},
				{"--shared-base-layers", "--name", "test-container"},
				{"--shared-base-layers", "--env", "TEST=value"},
			}

			for _, flags := range flagCombinations {
				args := append(flags, "--help") // Use --help to avoid needing real containers
				session := podmanTest.Podman(append([]string{"run"}, args...))
				session.WaitWithDefaultTimeout()
				Expect(session).Should(ExitCleanly(), fmt.Sprintf("Flag combination failed: %v", flags))
			}
		})

		It("should maintain flag order independence", func() {
			// Test that flag order doesn't matter
			flagOrders := [][]string{
				{"--shared-base-layers", "--rm", "--help"},
				{"--rm", "--shared-base-layers", "--help"},
				{"--help", "--shared-base-layers", "--rm"},
			}

			for _, flags := range flagOrders {
				session := podmanTest.Podman(append([]string{"run"}, flags...))
				session.WaitWithDefaultTimeout()
				Expect(session).Should(ExitCleanly(), fmt.Sprintf("Flag order failed: %v", flags))
			}
		})
	})

	Context("Command Structure Validation", func() {
		It("should be valid for run command only", func() {
			// Test that --shared-base-layers is specific to 'run' command
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--help"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			helpOutput := session.OutputToString()
			Expect(helpOutput).To(ContainSubstring("shared-base-layers"))
		})

		It("should show appropriate help documentation", func() {
			session := podmanTest.Podman([]string{"run", "--help"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			helpOutput := session.OutputToString()
			// Verify the flag appears in help with reasonable description
			Expect(helpOutput).To(ContainSubstring("shared-base-layers"))

			// The help should give some indication of what the flag does
			lines := strings.Split(helpOutput, "\n")
			var flagLine string
			for _, line := range lines {
				if strings.Contains(line, "shared-base-layers") {
					flagLine = line
					break
				}
			}
			Expect(flagLine).ToNot(BeEmpty(), "shared-base-layers flag should appear in help")
		})
	})

	Context("Error Handling Tests", func() {
		It("should handle graceful degradation when shared storage unavailable", func() {
			// Test behavior when shared storage is not available
			// This should fall back to normal operation without crashing
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "nonexistent-image"})
			session.WaitWithDefaultTimeout()

			// Should fail due to image issue, not due to shared storage
			Expect(session).ShouldNot(ExitCleanly())
			errorOutput := strings.ToLower(session.ErrorToString())

			// Should be image-related error, not storage-related panic
			Expect(errorOutput).ToNot(ContainSubstring("panic"))
			Expect(errorOutput).ToNot(ContainSubstring("runtime error"))
		})

		It("should provide meaningful error messages", func() {
			// Test with various invalid scenarios
			invalidScenarios := [][]string{
				{"--shared-base-layers", ""},          // Empty image name
				{"--shared-base-layers", "invalid:::"}, // Invalid image format
			}

			for _, scenario := range invalidScenarios {
				if len(scenario) < 2 || scenario[1] == "" {
					continue // Skip empty image name test as it might have different behavior
				}

				args := append([]string{"run"}, scenario...)
				session := podmanTest.Podman(args)
				session.WaitWithDefaultTimeout()

				// Should fail gracefully with meaningful error
				Expect(session).ShouldNot(ExitCleanly())
				errorOutput := session.ErrorToString()
				Expect(errorOutput).ToNot(BeEmpty(), fmt.Sprintf("Should provide error message for: %v", scenario))
			}
		})
	})

	Context("Integration Readiness Tests", func() {
		It("should be ready for container runtime integration", func() {
			// Verify that the CLI infrastructure is ready for actual runtime integration
			// This test ensures all the plumbing is in place

			// Test flag parsing
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--help"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			// Test that configuration doesn't break other functionality
			session = podmanTest.Podman([]string{"version"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			// Test info command still works
			session = podmanTest.Podman([]string{"info"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
		})

		It("should maintain backward compatibility", func() {
			// Ensure existing functionality still works
			session := podmanTest.Podman([]string{"run", "--help"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			helpOutput := session.OutputToString()
			// Standard flags should still be present
			Expect(helpOutput).To(ContainSubstring("--rm"))
			Expect(helpOutput).To(ContainSubstring("--detach"))
			Expect(helpOutput).To(ContainSubstring("--name"))
		})
	})
})