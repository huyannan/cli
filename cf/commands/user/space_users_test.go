package user_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/plugin/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/cli/testhelpers/matchers"
)

var _ = Describe("space-users command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		spaceRepo           *testapi.FakeSpaceRepository
		userRepo            *testapi.FakeUserRepository
		configRepo          core_config.Repository
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetUserRepository(userRepo)
		deps.RepoLocator = deps.RepoLocator.SetSpaceRepository(spaceRepo)

		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("space-users").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		configRepo = testconfig.NewRepositoryWithDefaults()
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		spaceRepo = &testapi.FakeSpaceRepository{}
		userRepo = &testapi.FakeUserRepository{}
		deps = command_registry.NewDependency()
		updateCommandDependency(false)
	})

	runCommand := func(args ...string) bool {
		cmd := command_registry.Commands.FindCommand("space-users")
		return testcmd.RunCliCommand(cmd, args, requirementsFactory)
	}

	Describe("requirements", func() {
		It("fails when not logged in", func() {
			Expect(runCommand("my-org", "my-space")).To(BeFalse())
		})

		It("succeeds when logged in", func() {
			requirementsFactory.LoginSuccess = true
			passed := runCommand("some-org", "whatever-space")

			Expect(passed).To(BeTrue())
			Expect("some-org").To(Equal(requirementsFactory.OrganizationName))
		})
	})

	It("fails with usage when not invoked with exactly two args", func() {
		runCommand("my-org")
		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Incorrect Usage", "Requires arguments"},
		))
	})

	Context("when logged in and given some users in the org and space", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true

			org := models.Organization{}
			org.Name = "Org1"
			org.Guid = "org1-guid"
			space := models.Space{}
			space.Name = "Space1"
			space.Guid = "space1-guid"

			requirementsFactory.Organization = org
			spaceRepo.FindByNameInOrgSpace = space

			user := models.UserFields{}
			user.Username = "user1"
			user2 := models.UserFields{}
			user2.Username = "user2"
			user3 := models.UserFields{}
			user3.Username = "user3"
			user4 := models.UserFields{}
			user4.Username = "user4"
			userRepo.ListUsersByRole = map[string][]models.UserFields{
				models.SPACE_MANAGER:   []models.UserFields{user, user2},
				models.SPACE_DEVELOPER: []models.UserFields{user4},
				models.SPACE_AUDITOR:   []models.UserFields{user3},
			}
		})

		It("tells you about the space users in the given space", func() {
			runCommand("my-org", "my-space")

			Expect(spaceRepo.FindByNameInOrgName).To(Equal("my-space"))
			Expect(spaceRepo.FindByNameInOrgOrgGuid).To(Equal("org1-guid"))
			Expect(userRepo.ListUsersSpaceGuid).To(Equal("space1-guid"))

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting users in org", "Org1", "Space1", "my-user"},
				[]string{"SPACE MANAGER"},
				[]string{"user1"},
				[]string{"user2"},
				[]string{"SPACE DEVELOPER"},
				[]string{"user4"},
				[]string{"SPACE AUDITOR"},
				[]string{"user3"},
			))
		})

		Context("when cc api verson is >= 2.21.0", func() {
			BeforeEach(func() {
				userRepo.ListUsersInSpaceForRole_CallCount = 0
				userRepo.ListUsersInSpaceForRoleWithNoUAA_CallCount = 0
			})

			It("calls ListUsersInSpaceForRoleWithNoUAA()", func() {
				configRepo.SetApiVersion("2.22.0")
				runCommand("my-org", "my-sapce")

				Expect(userRepo.ListUsersInSpaceForRoleWithNoUAA_CallCount).To(BeNumerically(">=", 1))
				Expect(userRepo.ListUsersInSpaceForRole_CallCount).To(Equal(0))
			})
		})

		Context("when cc api verson is < 2.21.0", func() {
			It("calls ListUsersInSpaceForRole()", func() {
				configRepo.SetApiVersion("2.20.0")
				runCommand("my-org", "my-space")

				Expect(userRepo.ListUsersInSpaceForRoleWithNoUAA_CallCount).To(Equal(0))
				Expect(userRepo.ListUsersInSpaceForRole_CallCount).To(BeNumerically(">=", 1))
			})
		})
	})

	Describe("when invoked by a plugin", func() {
		var (
			pluginUserModel []plugin_models.User
		)

		Context("single roles", func() {

			BeforeEach(func() {

				org := models.Organization{}
				org.Name = "the-org"
				org.Guid = "the-org-guid"

				space := models.Space{}
				space.Name = "the-space"
				space.Guid = "the-space-guid"

				// space managers
				user := models.UserFields{}
				user.Username = "user1"
				user.Guid = "1111"

				user2 := models.UserFields{}
				user2.Username = "user2"
				user2.Guid = "2222"

				// space auditor
				user3 := models.UserFields{}
				user3.Username = "user3"
				user3.Guid = "3333"

				// space developer
				user4 := models.UserFields{}
				user4.Username = "user4"
				user4.Guid = "4444"

				userRepo.ListUsersByRole = map[string][]models.UserFields{
					models.SPACE_MANAGER:   []models.UserFields{user, user2},
					models.SPACE_DEVELOPER: []models.UserFields{user4},
					models.SPACE_AUDITOR:   []models.UserFields{user3},
				}

				requirementsFactory.LoginSuccess = true
				requirementsFactory.Organization = org
				requirementsFactory.Space = space
				pluginUserModel = []plugin_models.User{}
				deps.PluginModels.Users = &pluginUserModel
				updateCommandDependency(true)
			})

			It("populates the plugin model with users with single roles", func() {
				runCommand("the-org", "the-space")
				Ω(len(pluginUserModel)).To(Equal(4))
				Ω(pluginUserModel[0].Username).To(Equal("user1"))
				Ω(pluginUserModel[0].Guid).To(Equal("1111"))
				Ω(pluginUserModel[0].Roles[0]).To(Equal(models.SPACE_MANAGER))

				Ω(pluginUserModel[1].Username).To(Equal("user2"))
				Ω(pluginUserModel[1].Guid).To(Equal("2222"))
				Ω(pluginUserModel[1].Roles[0]).To(Equal(models.SPACE_MANAGER))

				Ω(pluginUserModel[2].Username).To(Equal("user4"))
				Ω(pluginUserModel[2].Guid).To(Equal("4444"))
				Ω(pluginUserModel[2].Roles[0]).To(Equal(models.SPACE_DEVELOPER))

				Ω(pluginUserModel[3].Username).To(Equal("user3"))
				Ω(pluginUserModel[3].Guid).To(Equal("3333"))
				Ω(pluginUserModel[3].Roles[0]).To(Equal(models.SPACE_AUDITOR))

			})

		})

		Context("multiple roles", func() {

			BeforeEach(func() {

				org := models.Organization{}
				org.Name = "the-org"
				org.Guid = "the-org-guid"

				space := models.Space{}
				space.Name = "the-space"
				space.Guid = "the-space-guid"

				// space managers
				user := models.UserFields{}
				user.Username = "user1"
				user.Guid = "1111"

				user2 := models.UserFields{}
				user2.Username = "user2"
				user2.Guid = "2222"

				// space auditor
				user3 := models.UserFields{}
				user3.Username = "user3"
				user3.Guid = "3333"

				// space developer
				user4 := models.UserFields{}
				user4.Username = "user4"
				user4.Guid = "4444"

				userRepo.ListUsersByRole = map[string][]models.UserFields{
					models.SPACE_MANAGER:   []models.UserFields{user, user2, user3, user4},
					models.SPACE_DEVELOPER: []models.UserFields{user2, user4},
					models.SPACE_AUDITOR:   []models.UserFields{user, user3},
				}

				requirementsFactory.LoginSuccess = true
				requirementsFactory.Organization = org
				requirementsFactory.Space = space
				pluginUserModel = []plugin_models.User{}
				deps.PluginModels.Users = &pluginUserModel
				updateCommandDependency(true)
			})

			It("populates the plugin model with users with multiple roles", func() {
				runCommand("the-org", "the-space")
				Ω(len(pluginUserModel)).To(Equal(4))

				// user1
				Ω(len(pluginUserModel[0].Roles)).To(Equal(2))
				Ω(pluginUserModel[0].Roles[0]).To(Equal(models.SPACE_MANAGER))
				Ω(pluginUserModel[0].Roles[1]).To(Equal(models.SPACE_AUDITOR))

				// user2
				Ω(len(pluginUserModel[1].Roles)).To(Equal(2))
				Ω(pluginUserModel[1].Roles[0]).To(Equal(models.SPACE_MANAGER))
				Ω(pluginUserModel[1].Roles[1]).To(Equal(models.SPACE_DEVELOPER))

				// user3
				Ω(len(pluginUserModel[2].Roles)).To(Equal(2))
				Ω(pluginUserModel[2].Roles[0]).To(Equal(models.SPACE_MANAGER))
				Ω(pluginUserModel[2].Roles[1]).To(Equal(models.SPACE_AUDITOR))

				// user4
				Ω(len(pluginUserModel[3].Roles)).To(Equal(2))
				Ω(pluginUserModel[3].Roles[0]).To(Equal(models.SPACE_MANAGER))
				Ω(pluginUserModel[3].Roles[1]).To(Equal(models.SPACE_DEVELOPER))

			})

		})

	})

})
