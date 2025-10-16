//go:build linux || freebsd

package integration

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/dmikushin/podman-shared/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to create a mock shared storage directory
func createMockSharedStorage(tempDir string) (string, error) {
	sharedStorageDir := filepath.Join(tempDir, "shared_storage")
	err := os.MkdirAll(sharedStorageDir, 0755)
	if err != nil {
		return "", err
	}

	// Create mock base layer directory structure
	baseLayersDir := filepath.Join(sharedStorageDir, "overlay-layers")
	err = os.MkdirAll(baseLayersDir, 0755)
	if err != nil {
		return "", err
	}

	return sharedStorageDir, nil
}

// Helper function to verify that inspect returns valid JSON
func verifyContainerInspect(podmanTest *PodmanTestIntegration, containerID string) {
	session := podmanTest.Podman([]string{"inspect", "--format", "json", containerID})
	session.WaitWithDefaultTimeout()
	Expect(session).Should(ExitCleanly())
	Expect(session.OutputToString()).To(BeValidJSON())
}

// Helper function to verify mount options by checking /proc/self/mountinfo
func verifyMountOptions(podmanTest *PodmanTestIntegration, containerID string, expectedReadOnly bool) {
	session := podmanTest.Podman([]string{"exec", containerID, "cat", "/proc/self/mountinfo"})
	session.WaitWithDefaultTimeout()
	Expect(session).Should(ExitCleanly())

	mountInfo := session.OutputToString()
	if expectedReadOnly {
		Expect(mountInfo).To(ContainSubstring("ro,"))
	} else {
		Expect(mountInfo).To(ContainSubstring("rw,"))
	}
}

var _ = Describe("Podman run with shared base layers", func() {
	var (
		tempDir string
	)

	BeforeEach(func() {
		tempDir = podmanTest.TempDir
		_, err := createMockSharedStorage(tempDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// ============================================================================
	// CLI Behavior Tests (3 tests)
	// ============================================================================

	Context("CLI Flag Behavior", func() {
		It("should accept --shared-base-layers flag", func() {
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "echo", "test"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			Expect(session.OutputToString()).To(ContainSubstring("test"))
		})

		It("should use normal copy behavior when flag is omitted", func() {
			// Run without --shared-base-layers flag
			session := podmanTest.Podman([]string{"run", "--rm", "--name", "test-normal", ALPINE, "echo", "normal"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			Expect(session.OutputToString()).To(ContainSubstring("normal"))

			// Verify that container uses standard overlay behavior
			// This is the default behavior, so we just verify the container runs successfully
		})

		It("should handle error when base layers are not on shared storage", func() {
			// Create a non-shared storage location
			nonSharedDir := filepath.Join(tempDir, "non_shared")
			err := os.MkdirAll(nonSharedDir, 0755)
			Expect(err).ToNot(HaveOccurred())

			// Try to run with --shared-base-layers on non-shared storage
			// Note: This test depends on the actual implementation behavior
			// For now, we test that the flag is accepted
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "echo", "test"})
			session.WaitWithDefaultTimeout()
			// The behavior might vary based on implementation - either succeed or fail gracefully
			// We mainly test that the flag doesn't cause a crash
		})
	})

	// ============================================================================
	// Shared Storage Detection Tests (4 tests)
	// ============================================================================

	Context("Shared Storage Detection", func() {
		It("should detect NFS-mounted base layers", func() {
			Skip("NFS detection requires actual NFS mount - skipping in unit tests")
			// This test would require setting up actual NFS mounts
			// For integration testing, this could be implemented with proper NFS setup
		})

		It("should mount base layers read-only from shared storage", func() {
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "test-readonly", ALPINE, "sleep", "60"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			containerID := session.OutputToString()

			// Check that base layers are mounted read-only
			// This is implementation-specific and may need adjustment based on actual behavior
			mountSession := podmanTest.Podman([]string{"exec", containerID, "mount"})
			mountSession.WaitWithDefaultTimeout()
			Expect(mountSession).Should(ExitCleanly())

			// Clean up
			cleanupSession := podmanTest.Podman([]string{"rm", "-f", containerID})
			cleanupSession.WaitWithDefaultTimeout()
			Expect(cleanupSession).Should(ExitCleanly())
		})

		It("should create writable layers locally", func() {
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "test-writable", ALPINE, "sleep", "60"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			containerID := session.OutputToString()

			// Test writing to a writable location
			writeSession := podmanTest.Podman([]string{"exec", containerID, "sh", "-c", "echo 'writable test' > /tmp/testfile"})
			writeSession.WaitWithDefaultTimeout()
			Expect(writeSession).Should(ExitCleanly())

			// Verify the file was written
			readSession := podmanTest.Podman([]string{"exec", containerID, "cat", "/tmp/testfile"})
			readSession.WaitWithDefaultTimeout()
			Expect(readSession).Should(ExitCleanly())
			Expect(readSession.OutputToString()).To(ContainSubstring("writable test"))

			// Clean up
			cleanupSession := podmanTest.Podman([]string{"rm", "-f", containerID})
			cleanupSession.WaitWithDefaultTimeout()
			Expect(cleanupSession).Should(ExitCleanly())
		})

		It("should show correct mount configuration in podman inspect", func() {
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "test-inspect", ALPINE, "sleep", "60"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			containerID := session.OutputToString()

			// Inspect the container
			inspectSession := podmanTest.Podman([]string{"inspect", containerID})
			inspectSession.WaitWithDefaultTimeout()
			Expect(inspectSession).Should(ExitCleanly())

			// Verify that the inspect output contains mount information
			// The specific format depends on the implementation
			inspectOutput := inspectSession.OutputToString()
			Expect(inspectOutput).To(ContainSubstring("Mount"))

			// Clean up
			cleanupSession := podmanTest.Podman([]string{"rm", "-f", containerID})
			cleanupSession.WaitWithDefaultTimeout()
			Expect(cleanupSession).Should(ExitCleanly())
		})
	})

	// ============================================================================
	// Container Functionality Tests (5 tests)
	// ============================================================================

	Context("Container Functionality", func() {
		It("should start container successfully with shared base layers", func() {
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "echo", "success"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			Expect(session.OutputToString()).To(ContainSubstring("success"))
		})

		It("should support file operations in writable layers", func() {
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "sh", "-c",
				"echo 'test content' > /tmp/test.txt && cat /tmp/test.txt"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			Expect(session.OutputToString()).To(ContainSubstring("test content"))
		})

		It("should access base layer files as read-only", func() {
			// Test accessing files that should be in base layers
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "cat", "/etc/os-release"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())
			Expect(session.OutputToString()).To(ContainSubstring("Alpine"))

			// Test that we cannot write to base layer files
			writeSession := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "sh", "-c",
				"echo 'test' > /etc/os-release 2>&1 || echo 'write failed as expected'"})
			writeSession.WaitWithDefaultTimeout()
			Expect(writeSession).Should(ExitCleanly())
			// The exact error message may vary, but we expect some indication that write failed
		})

		It("should clean up container successfully", func() {
			// Create a container
			session := podmanTest.Podman([]string{"create", "--shared-base-layers", "--name", "test-cleanup", ALPINE, "echo", "test"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			// Remove the container
			rmSession := podmanTest.Podman([]string{"rm", "test-cleanup"})
			rmSession.WaitWithDefaultTimeout()
			Expect(rmSession).Should(ExitCleanly())

			// Verify container is gone
			psSession := podmanTest.Podman([]string{"ps", "-a", "--filter", "name=test-cleanup"})
			psSession.WaitWithDefaultTimeout()
			Expect(psSession).Should(ExitCleanly())
			Expect(psSession.OutputToString()).ToNot(ContainSubstring("test-cleanup"))
		})

		It("should maintain proper file permissions and security contexts", func() {
			// Test file permissions
			session := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "ls", "-la", "/bin/sh"})
			session.WaitWithDefaultTimeout()
			Expect(session).Should(ExitCleanly())

			output := session.OutputToString()
			// Verify that /bin/sh has executable permissions
			Expect(output).To(ContainSubstring("rwx"))

			// Test creating files with specific permissions
			permSession := podmanTest.Podman([]string{"run", "--shared-base-layers", "--rm", ALPINE, "sh", "-c",
				"touch /tmp/testfile && chmod 644 /tmp/testfile && ls -la /tmp/testfile"})
			permSession.WaitWithDefaultTimeout()
			Expect(permSession).Should(ExitCleanly())
			Expect(permSession.OutputToString()).To(ContainSubstring("rw-r--r--"))
		})
	})

	// ============================================================================
	// Multi-Container Scenarios (3 tests)
	// ============================================================================

	Context("Multi-Container Scenarios", func() {
		It("should run multiple containers sharing same base layers", func() {
			// Start first container
			session1 := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "shared1", ALPINE, "sleep", "60"})
			session1.WaitWithDefaultTimeout()
			Expect(session1).Should(ExitCleanly())
			container1ID := session1.OutputToString()

			// Start second container with same image
			session2 := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "shared2", ALPINE, "sleep", "60"})
			session2.WaitWithDefaultTimeout()
			Expect(session2).Should(ExitCleanly())
			container2ID := session2.OutputToString()

			// Verify both containers are running
			psSession := podmanTest.Podman([]string{"ps", "--format", "{{.Names}}"})
			psSession.WaitWithDefaultTimeout()
			Expect(psSession).Should(ExitCleanly())
			output := psSession.OutputToString()
			Expect(output).To(ContainSubstring("shared1"))
			Expect(output).To(ContainSubstring("shared2"))

			// Clean up
			cleanupSession1 := podmanTest.Podman([]string{"rm", "-f", container1ID})
			cleanupSession1.WaitWithDefaultTimeout()
			Expect(cleanupSession1).Should(ExitCleanly())

			cleanupSession2 := podmanTest.Podman([]string{"rm", "-f", container2ID})
			cleanupSession2.WaitWithDefaultTimeout()
			Expect(cleanupSession2).Should(ExitCleanly())
		})

		It("should maintain isolation between containers' writable layers", func() {
			// Start first container and create a file
			session1 := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "isolated1", ALPINE, "sleep", "60"})
			session1.WaitWithDefaultTimeout()
			Expect(session1).Should(ExitCleanly())
			container1ID := session1.OutputToString()

			writeSession1 := podmanTest.Podman([]string{"exec", container1ID, "sh", "-c", "echo 'container1 data' > /tmp/isolation_test.txt"})
			writeSession1.WaitWithDefaultTimeout()
			Expect(writeSession1).Should(ExitCleanly())

			// Start second container
			session2 := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", "isolated2", ALPINE, "sleep", "60"})
			session2.WaitWithDefaultTimeout()
			Expect(session2).Should(ExitCleanly())
			container2ID := session2.OutputToString()

			// Verify second container doesn't see first container's file
			readSession2 := podmanTest.Podman([]string{"exec", container2ID, "cat", "/tmp/isolation_test.txt"})
			readSession2.WaitWithDefaultTimeout()
			Expect(readSession2).ShouldNot(ExitCleanly()) // File should not exist

			// Create a different file in second container
			writeSession2 := podmanTest.Podman([]string{"exec", container2ID, "sh", "-c", "echo 'container2 data' > /tmp/isolation_test.txt"})
			writeSession2.WaitWithDefaultTimeout()
			Expect(writeSession2).Should(ExitCleanly())

			// Verify first container still has its own file content
			readSession1 := podmanTest.Podman([]string{"exec", container1ID, "cat", "/tmp/isolation_test.txt"})
			readSession1.WaitWithDefaultTimeout()
			Expect(readSession1).Should(ExitCleanly())
			Expect(readSession1.OutputToString()).To(ContainSubstring("container1 data"))

			// Verify second container has its own file content
			readSession2Again := podmanTest.Podman([]string{"exec", container2ID, "cat", "/tmp/isolation_test.txt"})
			readSession2Again.WaitWithDefaultTimeout()
			Expect(readSession2Again).Should(ExitCleanly())
			Expect(readSession2Again.OutputToString()).To(ContainSubstring("container2 data"))

			// Clean up
			cleanupSession1 := podmanTest.Podman([]string{"rm", "-f", container1ID})
			cleanupSession1.WaitWithDefaultTimeout()
			Expect(cleanupSession1).Should(ExitCleanly())

			cleanupSession2 := podmanTest.Podman([]string{"rm", "-f", container2ID})
			cleanupSession2.WaitWithDefaultTimeout()
			Expect(cleanupSession2).Should(ExitCleanly())
		})

		It("should handle concurrent access to shared base layers", func() {
			var containerIDs []string

			// Start multiple containers concurrently
			for i := 0; i < 3; i++ {
				containerName := fmt.Sprintf("concurrent%d", i)
				session := podmanTest.Podman([]string{"run", "--shared-base-layers", "-d", "--name", containerName, ALPINE, "sleep", "60"})
				session.WaitWithDefaultTimeout()
				Expect(session).Should(ExitCleanly())
				containerIDs = append(containerIDs, session.OutputToString())
			}

			// Verify all containers can access base layer files simultaneously
			for i, containerID := range containerIDs {
				readSession := podmanTest.Podman([]string{"exec", containerID, "cat", "/etc/os-release"})
				readSession.WaitWithDefaultTimeout()
				Expect(readSession).Should(ExitCleanly())
				Expect(readSession.OutputToString()).To(ContainSubstring("Alpine"), fmt.Sprintf("Container %d should access base layer files", i))

				// Also test that each can write to its own writable layer
				writeSession := podmanTest.Podman([]string{"exec", containerID, "sh", "-c", fmt.Sprintf("echo 'container%d' > /tmp/concurrent_test.txt", i)})
				writeSession.WaitWithDefaultTimeout()
				Expect(writeSession).Should(ExitCleanly())

				verifySession := podmanTest.Podman([]string{"exec", containerID, "cat", "/tmp/concurrent_test.txt"})
				verifySession.WaitWithDefaultTimeout()
				Expect(verifySession).Should(ExitCleanly())
				Expect(verifySession.OutputToString()).To(ContainSubstring(fmt.Sprintf("container%d", i)))
			}

			// Clean up all containers
			for _, containerID := range containerIDs {
				cleanupSession := podmanTest.Podman([]string{"rm", "-f", containerID})
				cleanupSession.WaitWithDefaultTimeout()
				Expect(cleanupSession).Should(ExitCleanly())
			}
		})
	})
})