/**
 * Generated by orval v6.10.3 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import {
  rest
} from 'msw'
import {
  faker
} from '@faker-js/faker'
import {
  AccessRuleStatus,
  RequestStatus,
  ApprovalMethod,
  IdpStatus
} from '.././types'

export const getAdminListAccessRulesMock = () => ({accessRules: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88d0003dycc9roobxy2': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88d0002dycc5m9ogsai': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})), next: faker.helpers.arrayElement([faker.random.word(), null])})

export const getAdminCreateAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88e0005dycc66ic2gox': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88e0004dyccat7u0n56': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})

export const getAdminGetAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88l0007dycc621s64lk': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88l0006dyccfu5hgfov': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})

export const getAdminUpdateAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88m0009dycce01w2g6z': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88m0008dyccfs6p79i1': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})

export const getAdminArchiveAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88r000bdyccaxe4eow0': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88r000adycc2zhkgr2p': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})

export const getAdminGetAccessRuleVersionsMock = () => ({accessRules: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88t000ddycc0871amdp': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88t000cdycc93bvdetf': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})), next: faker.helpers.arrayElement([faker.random.word(), null])})

export const getAdminGetAccessRuleVersionMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(Object.values(AccessRuleStatus)), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), approval: {users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.helpers.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cldl0w88z000fdycc61a05oie': {values: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), groupings: {
        'cldl0w88z000edyccbpvoc3ba': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))
      }, formElement: faker.helpers.arrayElement(['INPUT','MULTISELECT'])}
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number({min: 60, max: 15724800})}, isCurrent: faker.datatype.boolean()})

export const getAdminGetDeploymentVersionMock = () => ({version: faker.random.word()})

export const getAdminListRequestsMock = () => ({requests: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), requestor: faker.random.word(), status: faker.helpers.arrayElement(Object.values(RequestStatus)), reason: faker.helpers.arrayElement([faker.random.word(), undefined]), timing: {durationSeconds: faker.datatype.number({min: undefined, max: undefined}), startTime: faker.helpers.arrayElement([faker.random.word(), undefined])}, requestedAt: faker.random.word(), accessRuleId: faker.random.word(), accessRuleVersion: faker.random.word(), updatedAt: faker.random.word(), grant: faker.helpers.arrayElement([{status: faker.helpers.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), subject: faker.internet.email(), provider: faker.random.word(), start: `${faker.date.past().toISOString().split('.')[0]}Z`, end: `${faker.date.past().toISOString().split('.')[0]}Z`}, undefined]), approvalMethod: faker.helpers.arrayElement([faker.helpers.arrayElement(Object.values(ApprovalMethod)), undefined])})), next: faker.helpers.arrayElement([faker.random.word(), null])})

export const getAdminUpdateUserMock = () => ({id: faker.random.word(), email: faker.random.word(), firstName: faker.random.word(), picture: faker.random.word(), status: faker.helpers.arrayElement(Object.values(IdpStatus)), lastName: faker.random.word(), updatedAt: faker.random.word(), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))})

export const getAdminListUsersMock = () => ({users: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), email: faker.random.word(), firstName: faker.random.word(), picture: faker.random.word(), status: faker.helpers.arrayElement(Object.values(IdpStatus)), lastName: faker.random.word(), updatedAt: faker.random.word(), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))})), next: faker.helpers.arrayElement([faker.random.word(), null])})

export const getAdminCreateUserMock = () => ({id: faker.random.word(), email: faker.random.word(), firstName: faker.random.word(), picture: faker.random.word(), status: faker.helpers.arrayElement(Object.values(IdpStatus)), lastName: faker.random.word(), updatedAt: faker.random.word(), groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word()))})

export const getAdminListGroupsMock = () => ({groups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({name: faker.random.word(), description: faker.random.word(), id: faker.random.word(), memberCount: faker.datatype.number({min: undefined, max: undefined}), members: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), source: faker.random.word()})), next: faker.helpers.arrayElement([faker.random.word(), null])})

export const getAdminCreateGroupMock = () => ({name: faker.random.word(), description: faker.random.word(), id: faker.random.word(), memberCount: faker.datatype.number({min: undefined, max: undefined}), members: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), source: faker.random.word()})

export const getAdminGetGroupMock = () => ({name: faker.random.word(), description: faker.random.word(), id: faker.random.word(), memberCount: faker.datatype.number({min: undefined, max: undefined}), members: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), source: faker.random.word()})

export const getAdminUpdateGroupMock = () => ({name: faker.random.word(), description: faker.random.word(), id: faker.random.word(), memberCount: faker.datatype.number({min: undefined, max: undefined}), members: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), source: faker.random.word()})

export const getAdminListProvidersMock = () => (Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), type: faker.random.word()})))

export const getAdminListProvidersv2Mock = () => (Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), type: faker.random.word()})))

export const getAdminGetProviderv2Mock = () => ({id: faker.random.word(), type: faker.random.word()})

export const getAdminGetProviderMock = () => ({id: faker.random.word(), type: faker.random.word()})

export const getAdminGetProviderArgsMock = () => ({
        'cldl0w89z000idycc534q75l7': {id: faker.random.word(), title: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined]), ruleFormElement: faker.helpers.arrayElement(['INPUT','MULTISELECT','SELECT']), requestFormElement: faker.helpers.arrayElement([faker.helpers.arrayElement(['SELECT']), undefined]), groups: faker.helpers.arrayElement([{
        'cldl0w89z000hdycc2vcz0v21': {id: faker.random.word(), title: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined])}
      }, undefined])}
      })

export const getAdminListProviderArgOptionsMock = () => ({options: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({label: faker.random.word(), value: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined])})), groups: faker.helpers.arrayElement([{
        'cldl0w8a0000jdyccd07yb3gm': Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({label: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined]), value: faker.random.word(), children: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), labelPrefix: faker.helpers.arrayElement([faker.random.word(), undefined])}))
      }, undefined])})

export const getAdminListProvidersetupsMock = () => ({providerSetups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), type: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(['COMPLETE','VALIDATION_FAILED','VALIDATING','INITIAL_CONFIGURATION_IN_PROGRESS','VALIDATION_SUCEEDED']), steps: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({complete: faker.datatype.boolean()})), configValues: {}, configValidation: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), status: faker.helpers.arrayElement(['IN_PROGRESS','SUCCESS','PENDING','ERROR']), fieldsValidated: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), logs: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.helpers.arrayElement(['INFO','WARNING','ERROR']), msg: faker.random.word()}))}))}))})

export const getAdminCreateProvidersetupMock = () => ({id: faker.random.word(), type: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(['COMPLETE','VALIDATION_FAILED','VALIDATING','INITIAL_CONFIGURATION_IN_PROGRESS','VALIDATION_SUCEEDED']), steps: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({complete: faker.datatype.boolean()})), configValues: {}, configValidation: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), status: faker.helpers.arrayElement(['IN_PROGRESS','SUCCESS','PENDING','ERROR']), fieldsValidated: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), logs: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.helpers.arrayElement(['INFO','WARNING','ERROR']), msg: faker.random.word()}))}))})

export const getAdminGetProvidersetupMock = () => ({id: faker.random.word(), type: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(['COMPLETE','VALIDATION_FAILED','VALIDATING','INITIAL_CONFIGURATION_IN_PROGRESS','VALIDATION_SUCEEDED']), steps: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({complete: faker.datatype.boolean()})), configValues: {}, configValidation: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), status: faker.helpers.arrayElement(['IN_PROGRESS','SUCCESS','PENDING','ERROR']), fieldsValidated: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), logs: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.helpers.arrayElement(['INFO','WARNING','ERROR']), msg: faker.random.word()}))}))})

export const getAdminDeleteProvidersetupMock = () => ({id: faker.random.word(), type: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(['COMPLETE','VALIDATION_FAILED','VALIDATING','INITIAL_CONFIGURATION_IN_PROGRESS','VALIDATION_SUCEEDED']), steps: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({complete: faker.datatype.boolean()})), configValues: {}, configValidation: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), status: faker.helpers.arrayElement(['IN_PROGRESS','SUCCESS','PENDING','ERROR']), fieldsValidated: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), logs: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.helpers.arrayElement(['INFO','WARNING','ERROR']), msg: faker.random.word()}))}))})

export const getAdminGetProvidersetupInstructionsMock = () => ({stepDetails: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({title: faker.random.word(), instructions: faker.random.word(), configFields: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), description: faker.random.word(), isOptional: faker.datatype.boolean(), isSecret: faker.datatype.boolean(), secretPath: faker.helpers.arrayElement([faker.random.word(), undefined])}))}))})

export const getAdminValidateProvidersetupMock = () => ({id: faker.random.word(), type: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(['COMPLETE','VALIDATION_FAILED','VALIDATING','INITIAL_CONFIGURATION_IN_PROGRESS','VALIDATION_SUCEEDED']), steps: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({complete: faker.datatype.boolean()})), configValues: {}, configValidation: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), status: faker.helpers.arrayElement(['IN_PROGRESS','SUCCESS','PENDING','ERROR']), fieldsValidated: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), logs: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.helpers.arrayElement(['INFO','WARNING','ERROR']), msg: faker.random.word()}))}))})

export const getAdminCompleteProvidersetupMock = () => ({deploymentConfigUpdateRequired: faker.datatype.boolean()})

export const getAdminSubmitProvidersetupStepMock = () => ({id: faker.random.word(), type: faker.random.word(), version: faker.random.word(), status: faker.helpers.arrayElement(['COMPLETE','VALIDATION_FAILED','VALIDATING','INITIAL_CONFIGURATION_IN_PROGRESS','VALIDATION_SUCEEDED']), steps: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({complete: faker.datatype.boolean()})), configValues: {}, configValidation: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), name: faker.random.word(), status: faker.helpers.arrayElement(['IN_PROGRESS','SUCCESS','PENDING','ERROR']), fieldsValidated: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => (faker.random.word())), logs: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.helpers.arrayElement(['INFO','WARNING','ERROR']), msg: faker.random.word()}))}))})

export const getAdminGetIdentityConfigurationMock = () => ({identityProvider: faker.random.word(), administratorGroupId: faker.random.word()})

export const getAdminMSW = () => [
rest.get('*/api/v1/admin/access-rules', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListAccessRulesMock()),
        )
      }),rest.post('*/api/v1/admin/access-rules', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminCreateAccessRuleMock()),
        )
      }),rest.get('*/api/v1/admin/access-rules/:ruleId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetAccessRuleMock()),
        )
      }),rest.put('*/api/v1/admin/access-rules/:ruleId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminUpdateAccessRuleMock()),
        )
      }),rest.post('*/api/v1/admin/access-rules/:ruleId/archive', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminArchiveAccessRuleMock()),
        )
      }),rest.get('*/api/v1/admin/access-rules/:ruleId/versions', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetAccessRuleVersionsMock()),
        )
      }),rest.get('*/api/v1/admin/access-rules/:ruleId/versions/:version', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetAccessRuleVersionMock()),
        )
      }),rest.get('*/api/v1/admin/deployment/version', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetDeploymentVersionMock()),
        )
      }),rest.get('*/api/v1/admin/requests', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListRequestsMock()),
        )
      }),rest.post('*/api/v1/admin/users/:userId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminUpdateUserMock()),
        )
      }),rest.get('*/api/v1/admin/users', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListUsersMock()),
        )
      }),rest.post('*/api/v1/admin/users', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminCreateUserMock()),
        )
      }),rest.get('*/api/v1/admin/groups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListGroupsMock()),
        )
      }),rest.post('*/api/v1/admin/groups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminCreateGroupMock()),
        )
      }),rest.get('*/api/v1/admin/groups/:groupId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetGroupMock()),
        )
      }),rest.put('*/api/v1/admin/groups/:groupId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminUpdateGroupMock()),
        )
      }),rest.delete('*/api/v1/admin/groups/:groupId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.get('*/api/v1/admin/providers', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListProvidersMock()),
        )
      }),rest.get('*/api/v1/admin/providersv2', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListProvidersv2Mock()),
        )
      }),rest.get('*/api/v1/admin/providersv2/:providerId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetProviderv2Mock()),
        )
      }),rest.get('*/api/v1/admin/providers/:providerId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetProviderMock()),
        )
      }),rest.get('*/api/v1/admin/providers/:providerId/args', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetProviderArgsMock()),
        )
      }),rest.get('*/api/v1/admin/providers/:providerId/args/:argId/options', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListProviderArgOptionsMock()),
        )
      }),rest.get('*/api/v1/admin/providersetups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListProvidersetupsMock()),
        )
      }),rest.post('*/api/v1/admin/providersetups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminCreateProvidersetupMock()),
        )
      }),rest.get('*/api/v1/admin/providersetups/:providersetupId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetProvidersetupMock()),
        )
      }),rest.delete('*/api/v1/admin/providersetups/:providersetupId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminDeleteProvidersetupMock()),
        )
      }),rest.get('*/api/v1/admin/providersetups/:providersetupId/instructions', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetProvidersetupInstructionsMock()),
        )
      }),rest.post('*/api/v1/admin/providersetups/:providersetupId/validate', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminValidateProvidersetupMock()),
        )
      }),rest.post('*/api/v1/admin/providersetups/:providersetupId/complete', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminCompleteProvidersetupMock()),
        )
      }),rest.put('*/api/v1/admin/providersetups/:providersetupId/steps/:stepIndex/complete', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminSubmitProvidersetupStepMock()),
        )
      }),rest.post('*/api/v1/admin/identity/sync', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.get('*/api/v1/admin/identity', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetIdentityConfigurationMock()),
        )
      }),]
