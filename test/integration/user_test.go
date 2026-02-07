package integration

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

var _ = Describe("User", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a User with all fields", func() {
		It("Should persist the User resource with name, email, and githubToken", func() {
			By("Creating a namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-user-all-fields",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			By("Creating a Secret with GITHUB_TOKEN")
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-github-token",
					Namespace: ns.Name,
				},
				StringData: map[string]string{
					"GITHUB_TOKEN": "test-gh-token",
				},
			}
			Expect(k8sClient.Create(ctx, secret)).Should(Succeed())

			By("Creating a User")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name:  "Test User",
					Email: "test@example.com",
					GitHubToken: &axonv1alpha1.SecretReference{
						Name: "user-github-token",
					},
				},
			}
			Expect(k8sClient.Create(ctx, user)).Should(Succeed())

			By("Verifying the User is created")
			userLookupKey := types.NamespacedName{Name: user.Name, Namespace: ns.Name}
			createdUser := &axonv1alpha1.User{}

			Eventually(func() error {
				return k8sClient.Get(ctx, userLookupKey, createdUser)
			}, timeout, interval).Should(Succeed())

			By("Verifying the User spec")
			Expect(createdUser.Spec.Name).To(Equal("Test User"))
			Expect(createdUser.Spec.Email).To(Equal("test@example.com"))
			Expect(createdUser.Spec.GitHubToken).NotTo(BeNil())
			Expect(createdUser.Spec.GitHubToken.Name).To(Equal("user-github-token"))

			By("Deleting the User")
			Expect(k8sClient.Delete(ctx, createdUser)).Should(Succeed())

			By("Verifying the User is deleted")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, userLookupKey, createdUser)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("When creating a User with only required fields", func() {
		It("Should persist the User resource with just the name", func() {
			By("Creating a namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-user-required-only",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			By("Creating a User with only name")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "minimal-user",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name: "Minimal User",
				},
			}
			Expect(k8sClient.Create(ctx, user)).Should(Succeed())

			By("Verifying the User is created")
			userLookupKey := types.NamespacedName{Name: user.Name, Namespace: ns.Name}
			createdUser := &axonv1alpha1.User{}

			Eventually(func() error {
				return k8sClient.Get(ctx, userLookupKey, createdUser)
			}, timeout, interval).Should(Succeed())

			By("Verifying the User spec")
			Expect(createdUser.Spec.Name).To(Equal("Minimal User"))
			Expect(createdUser.Spec.Email).To(BeEmpty())
			Expect(createdUser.Spec.GitHubToken).To(BeNil())
		})
	})

	Context("When creating a User with email but no github token", func() {
		It("Should persist the User resource with name and email", func() {
			By("Creating a namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-user-name-email",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			By("Creating a User with name and email")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "email-user",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name:  "Email User",
					Email: "email@example.com",
				},
			}
			Expect(k8sClient.Create(ctx, user)).Should(Succeed())

			By("Verifying the User is created")
			userLookupKey := types.NamespacedName{Name: user.Name, Namespace: ns.Name}
			createdUser := &axonv1alpha1.User{}

			Eventually(func() error {
				return k8sClient.Get(ctx, userLookupKey, createdUser)
			}, timeout, interval).Should(Succeed())

			By("Verifying the User spec")
			Expect(createdUser.Spec.Name).To(Equal("Email User"))
			Expect(createdUser.Spec.Email).To(Equal("email@example.com"))
			Expect(createdUser.Spec.GitHubToken).To(BeNil())
		})
	})

	Context("When updating a User", func() {
		It("Should update the User resource", func() {
			By("Creating a namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-user-update",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			By("Creating a User")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "update-user",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name:  "Original Name",
					Email: "original@example.com",
				},
			}
			Expect(k8sClient.Create(ctx, user)).Should(Succeed())

			By("Updating the User")
			userLookupKey := types.NamespacedName{Name: user.Name, Namespace: ns.Name}
			createdUser := &axonv1alpha1.User{}

			Eventually(func() error {
				return k8sClient.Get(ctx, userLookupKey, createdUser)
			}, timeout, interval).Should(Succeed())

			createdUser.Spec.Name = "Updated Name"
			createdUser.Spec.Email = "updated@example.com"
			Expect(k8sClient.Update(ctx, createdUser)).Should(Succeed())

			By("Verifying the User is updated")
			updatedUser := &axonv1alpha1.User{}
			Eventually(func() string {
				err := k8sClient.Get(ctx, userLookupKey, updatedUser)
				if err != nil {
					return ""
				}
				return updatedUser.Spec.Name
			}, timeout, interval).Should(Equal("Updated Name"))
			Expect(updatedUser.Spec.Email).To(Equal("updated@example.com"))
		})
	})

	Context("When listing Users", func() {
		It("Should list all Users in a namespace", func() {
			By("Creating a namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-user-list",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			By("Creating multiple Users")
			for _, name := range []string{"user-a", "user-b", "user-c"} {
				user := &axonv1alpha1.User{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns.Name,
					},
					Spec: axonv1alpha1.UserSpec{
						Name: name,
					},
				}
				Expect(k8sClient.Create(ctx, user)).Should(Succeed())
			}

			By("Listing Users in the namespace")
			userList := &axonv1alpha1.UserList{}
			Eventually(func() int {
				err := k8sClient.List(ctx, userList, &client.ListOptions{
					Namespace: ns.Name,
				})
				if err != nil {
					return 0
				}
				return len(userList.Items)
			}, timeout, interval).Should(Equal(3))
		})
	})
})
