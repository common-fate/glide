/**
 * Generated by orval v6.8.1 🍺
 * Do not edit manually.
 * Approvals
 * Granted Approvals API
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

export const getAdminListAccessRulesMock = () => ({accessRules: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), version: faker.random.word(), status: faker.random.arrayElement(Object.values(AccessRuleStatus)), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), approval: {users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.random.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8ox0003h8xd6djgb713': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getAdminCreateAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.random.arrayElement(Object.values(AccessRuleStatus)), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), approval: {users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.random.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8oy0004h8xdbvsadq3h': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})

export const getAdminGetAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.random.arrayElement(Object.values(AccessRuleStatus)), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), approval: {users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.random.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8p20005h8xd37wz8ask': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})

export const getAdminUpdateAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.random.arrayElement(Object.values(AccessRuleStatus)), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), approval: {users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.random.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8p20006h8xdhmeo425w': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})

export const getAdminGetAccessRuleVersionsMock = () => ({accessRules: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), version: faker.random.word(), status: faker.random.arrayElement(Object.values(AccessRuleStatus)), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), approval: {users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.random.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8p60008h8xd3pcm98n0': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getAdminGetAccessRuleVersionMock = () => ({id: faker.random.word(), version: faker.random.word(), status: faker.random.arrayElement(Object.values(AccessRuleStatus)), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), approval: {users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word()))}, name: faker.random.word(), description: faker.random.word(), metadata: {createdAt: faker.random.word(), createdBy: faker.random.word(), updatedAt: faker.random.word(), updatedBy: faker.random.word(), updateMessage: faker.random.arrayElement([faker.random.word(), undefined])}, target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8p80009h8xdb4mabsgw': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})

export const getAdminListRequestsMock = () => ({requests: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestor: faker.random.word(), status: faker.random.arrayElement(Object.values(RequestStatus)), reason: faker.random.arrayElement([faker.random.word(), undefined]), timing: {durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, requestedAt: faker.random.word(), accessRule: {id: faker.random.word(), version: faker.random.word()}, updatedAt: faker.random.word(), grant: faker.random.arrayElement([{status: faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), subject: faker.internet.email(), provider: faker.random.word(), start: faker.random.word(), end: faker.random.word()}, undefined]), approvalMethod: faker.random.arrayElement([faker.random.arrayElement(Object.values(ApprovalMethod)), undefined])})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getGetUsersMock = () => ({users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), email: faker.random.word(), firstName: faker.random.word(), picture: faker.random.word(), status: faker.random.arrayElement(Object.values(IdpStatus)), lastName: faker.random.word(), updatedAt: faker.random.word()})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getGetGroupsMock = () => ({groups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({name: faker.random.word(), description: faker.random.word(), id: faker.random.word()})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getGetGroupMock = () => ({name: faker.random.word(), description: faker.random.word(), id: faker.random.word()})

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
      }),rest.get('*/api/v1/admin/requests', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListRequestsMock()),
        )
      }),rest.get('*/api/v1/admin/users', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetUsersMock()),
        )
      }),rest.get('*/api/v1/admin/groups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetGroupsMock()),
        )
      }),rest.get('*/api/v1/admin/groups/:groupId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetGroupMock()),
        )
      }),]
