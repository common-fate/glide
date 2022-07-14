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
  RequestStatus,
  ApprovalMethod,
  IdpStatus
} from '.././types'

export const getListUserAccessRulesMock = () => ({accessRules: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), version: faker.random.word(), name: faker.random.word(), description: faker.random.word(), target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8n50000h8xd02nke166': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getUserGetAccessRuleMock = () => ({id: faker.random.word(), version: faker.random.word(), name: faker.random.word(), description: faker.random.word(), target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8nc0001h8xd5gmx4i0b': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()})

export const getUserGetAccessRuleApproversMock = () => ({users: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), next: faker.random.arrayElement([faker.random.word(), null])})

export const getUserListRequestsMock = () => ({requests: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestor: faker.random.word(), status: faker.random.arrayElement(Object.values(RequestStatus)), reason: faker.random.arrayElement([faker.random.word(), undefined]), timing: {durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, requestedAt: faker.random.word(), accessRule: {id: faker.random.word(), version: faker.random.word()}, updatedAt: faker.random.word(), grant: faker.random.arrayElement([{status: faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), subject: faker.internet.email(), provider: faker.random.word(), start: faker.random.word(), end: faker.random.word()}, undefined]), approvalMethod: faker.random.arrayElement([faker.random.arrayElement(Object.values(ApprovalMethod)), undefined])})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getUserGetRequestMock = () => ({id: faker.random.word(), requestor: faker.random.word(), status: faker.random.arrayElement(Object.values(RequestStatus)), reason: faker.random.arrayElement([faker.random.word(), undefined]), timing: {durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, requestedAt: faker.random.word(), accessRule: {id: faker.random.word(), version: faker.random.word(), name: faker.random.word(), description: faker.random.word(), target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8o70002h8xdbakncnry': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()}, updatedAt: faker.random.word(), grant: faker.random.arrayElement([{status: faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), subject: faker.internet.email(), provider: faker.random.word(), start: faker.random.word(), end: faker.random.word()}, undefined]), canReview: faker.datatype.boolean(), approvalMethod: faker.random.arrayElement([faker.random.arrayElement(Object.values(ApprovalMethod)), undefined])})

export const getListRequestEventsMock = () => ({events: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), createdAt: faker.random.word(), actor: faker.random.arrayElement([faker.random.word(), undefined]), fromStatus: faker.random.arrayElement([faker.random.arrayElement(Object.values(RequestStatus)), undefined]), toStatus: faker.random.arrayElement([faker.random.arrayElement(Object.values(RequestStatus)), undefined]), fromTiming: faker.random.arrayElement([{durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, undefined]), toTiming: faker.random.arrayElement([{durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, undefined]), fromGrantStatus: faker.random.arrayElement([faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), undefined]), toGrantStatus: faker.random.arrayElement([faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), undefined]), grantCreated: faker.random.arrayElement([faker.datatype.boolean(), undefined]), requestCreated: faker.random.arrayElement([faker.datatype.boolean(), undefined])})), next: faker.random.arrayElement([faker.random.word(), null])})

export const getReviewRequestMock = () => ({request: faker.random.arrayElement([{id: faker.random.word(), requestor: faker.random.word(), status: faker.random.arrayElement(Object.values(RequestStatus)), reason: faker.random.arrayElement([faker.random.word(), undefined]), timing: {durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, requestedAt: faker.random.word(), accessRule: {id: faker.random.word(), version: faker.random.word()}, updatedAt: faker.random.word(), grant: faker.random.arrayElement([{status: faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), subject: faker.internet.email(), provider: faker.random.word(), start: faker.random.word(), end: faker.random.word()}, undefined]), approvalMethod: faker.random.arrayElement([faker.random.arrayElement(Object.values(ApprovalMethod)), undefined])}, undefined])})

export const getCancelRequestMock = () => ({})

export const getGetAccessInstructionsMock = () => ({instructions: faker.random.arrayElement([faker.random.word(), undefined])})

export const getGetUserMock = () => ({id: faker.random.word(), email: faker.random.word(), firstName: faker.random.word(), picture: faker.random.word(), status: faker.random.arrayElement(Object.values(IdpStatus)), lastName: faker.random.word(), updatedAt: faker.random.word()})

export const getGetMeMock = () => ({user: {id: faker.random.word(), email: faker.random.word(), firstName: faker.random.word(), picture: faker.random.word(), status: faker.random.arrayElement(Object.values(IdpStatus)), lastName: faker.random.word(), updatedAt: faker.random.word()}, isAdmin: faker.datatype.boolean()})

export const getAdminGetRequestMock = () => ({id: faker.random.word(), requestor: faker.random.word(), status: faker.random.arrayElement(Object.values(RequestStatus)), reason: faker.random.arrayElement([faker.random.word(), undefined]), timing: {durationSeconds: faker.datatype.number(), startTime: faker.random.arrayElement([faker.random.word(), undefined])}, requestedAt: faker.random.word(), accessRule: {id: faker.random.word(), version: faker.random.word(), name: faker.random.word(), description: faker.random.word(), target: {provider: {id: faker.random.word(), type: faker.random.word()}, with: {
        'cl5kry8pc000ah8xd75vz1vm7': faker.random.word()
      }}, timeConstraints: {maxDurationSeconds: faker.datatype.number()}, isCurrent: faker.datatype.boolean()}, updatedAt: faker.random.word(), grant: faker.random.arrayElement([{status: faker.random.arrayElement(['PENDING','ACTIVE','ERROR','REVOKED','EXPIRED']), subject: faker.internet.email(), provider: faker.random.word(), start: faker.random.word(), end: faker.random.word()}, undefined]), canReview: faker.datatype.boolean(), approvalMethod: faker.random.arrayElement([faker.random.arrayElement(Object.values(ApprovalMethod)), undefined])})

export const getEndUserMSW = () => [
rest.get('*/api/v1/access-rules', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getListUserAccessRulesMock()),
        )
      }),rest.get('*/api/v1/access-rules/:ruleId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserGetAccessRuleMock()),
        )
      }),rest.get('*/api/v1/access-rules/:ruleId/approvers', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserGetAccessRuleApproversMock()),
        )
      }),rest.get('*/api/v1/requests', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserListRequestsMock()),
        )
      }),rest.post('*/api/v1/requests', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.get('*/api/v1/requests/:requestId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserGetRequestMock()),
        )
      }),rest.get('*/api/v1/requests/:requestId/events', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getListRequestEventsMock()),
        )
      }),rest.post('*/api/v1/requests/:requestId/review', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getReviewRequestMock()),
        )
      }),rest.post('*/api/v1/requests/:requestId/cancel', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getCancelRequestMock()),
        )
      }),rest.post('*/api/v1/requests/:requestid/revoke', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.get('*/api/v1/requests/:requestId/access-instructions', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetAccessInstructionsMock()),
        )
      }),rest.get('*/api/v1/users/:userId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetUserMock()),
        )
      }),rest.get('*/api/v1/users/me', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetMeMock()),
        )
      }),rest.get('*/api/v1/admin/requests/:requestId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminGetRequestMock()),
        )
      }),]
