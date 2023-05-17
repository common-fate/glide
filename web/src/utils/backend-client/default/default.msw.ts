/**
 * Generated by orval v6.7.1 🍺
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
  LogLevel,
  RequestStatus,
  RequestAccessGroupStatus,
  RequestAccessGroupApprovalMethod,
  RequestAccessGroupTargetStatus
} from '.././types'

export const getAdminFilterTargetGroupResourcesMock = () => ([...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({targetGroupId: faker.random.word(), resourceType: faker.random.word(), resource: {id: faker.random.word(), name: faker.random.word(), attributes: {}}})))

export const getAdminListTargetRoutesMock = () => ({routes: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({targetGroupId: faker.random.word(), handlerId: faker.random.word(), kind: faker.random.word(), priority: faker.datatype.number(), valid: faker.datatype.boolean(), diagnostics: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({level: faker.helpers.randomize(Object.values(LogLevel)), code: faker.random.word(), message: faker.random.word()}))})), next: faker.helpers.randomize([faker.random.word(), undefined])})

export const getUserListEntitlementsMock = () => ({entitlements: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}))})

export const getUserListEntitlementTargetsMock = () => ({targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), kind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()}))})), next: faker.helpers.randomize([faker.random.word(), undefined])})

export const getUserGetPreflightMock = () => ({id: faker.random.word(), accessGroups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requiresApproval: faker.datatype.boolean(), timeConstraints: {maxDurationSeconds: faker.datatype.number()}, targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), kind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()}))}))})), createdAt: faker.random.word()})

export const getUserRequestPreflightMock = () => ({id: faker.random.word(), accessGroups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requiresApproval: faker.datatype.boolean(), timeConstraints: {maxDurationSeconds: faker.datatype.number()}, targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), kind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()}))}))})), createdAt: faker.random.word()})

export const getUserListReviewsMock = () => ({requests: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), purpose: {reason: faker.helpers.randomize([faker.random.word(), undefined])}, accessGroups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestAccessGroupStatus)), requestedTiming: {durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, overrideTiming: faker.helpers.randomize([{durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, undefined]), updatedAt: faker.random.word(), createdAt: faker.random.word(), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), accessGroupId: faker.random.word(), targetGroupId: faker.random.word(), targetKind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()})), status: faker.helpers.randomize(Object.values(RequestAccessGroupTargetStatus)), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}})), approvalMethod: faker.helpers.randomize([faker.helpers.randomize(Object.values(RequestAccessGroupApprovalMethod)), undefined]), accessRule: {timeConstraints: {maxDurationSeconds: faker.datatype.number()}}, requestStatus: faker.helpers.randomize(Object.values(RequestStatus)), requestReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), groupReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), finalTiming: faker.helpers.randomize([{startTime: faker.random.word(), endTime: faker.random.word()}, undefined])})), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, requestedAt: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestStatus)), targetCount: faker.datatype.number()})), next: faker.helpers.randomize([faker.random.word(), null])})

export const getUserListRequestsMock = () => ({requests: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), purpose: {reason: faker.helpers.randomize([faker.random.word(), undefined])}, accessGroups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestAccessGroupStatus)), requestedTiming: {durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, overrideTiming: faker.helpers.randomize([{durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, undefined]), updatedAt: faker.random.word(), createdAt: faker.random.word(), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), accessGroupId: faker.random.word(), targetGroupId: faker.random.word(), targetKind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()})), status: faker.helpers.randomize(Object.values(RequestAccessGroupTargetStatus)), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}})), approvalMethod: faker.helpers.randomize([faker.helpers.randomize(Object.values(RequestAccessGroupApprovalMethod)), undefined]), accessRule: {timeConstraints: {maxDurationSeconds: faker.datatype.number()}}, requestStatus: faker.helpers.randomize(Object.values(RequestStatus)), requestReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), groupReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), finalTiming: faker.helpers.randomize([{startTime: faker.random.word(), endTime: faker.random.word()}, undefined])})), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, requestedAt: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestStatus)), targetCount: faker.datatype.number()})), next: faker.helpers.randomize([faker.random.word(), null])})

export const getUserPostRequestsMock = () => ({id: faker.random.word(), purpose: {reason: faker.helpers.randomize([faker.random.word(), undefined])}, accessGroups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestAccessGroupStatus)), requestedTiming: {durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, overrideTiming: faker.helpers.randomize([{durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, undefined]), updatedAt: faker.random.word(), createdAt: faker.random.word(), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), accessGroupId: faker.random.word(), targetGroupId: faker.random.word(), targetKind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()})), status: faker.helpers.randomize(Object.values(RequestAccessGroupTargetStatus)), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}})), approvalMethod: faker.helpers.randomize([faker.helpers.randomize(Object.values(RequestAccessGroupApprovalMethod)), undefined]), accessRule: {timeConstraints: {maxDurationSeconds: faker.datatype.number()}}, requestStatus: faker.helpers.randomize(Object.values(RequestStatus)), requestReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), groupReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), finalTiming: faker.helpers.randomize([{startTime: faker.random.word(), endTime: faker.random.word()}, undefined])})), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, requestedAt: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestStatus)), targetCount: faker.datatype.number()})

export const getUserGetRequestMock = () => ({id: faker.random.word(), purpose: {reason: faker.helpers.randomize([faker.random.word(), undefined])}, accessGroups: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestAccessGroupStatus)), requestedTiming: {durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, overrideTiming: faker.helpers.randomize([{durationSeconds: faker.datatype.number(), startTime: faker.helpers.randomize([faker.random.word(), undefined])}, undefined]), updatedAt: faker.random.word(), createdAt: faker.random.word(), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, targets: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), requestId: faker.random.word(), accessGroupId: faker.random.word(), targetGroupId: faker.random.word(), targetKind: {publisher: faker.random.word(), name: faker.random.word(), kind: faker.random.word(), icon: faker.random.word()}, fields: [...Array(faker.datatype.number({min: 1, max: 10}))].map(() => ({id: faker.random.word(), fieldTitle: faker.random.word(), fieldDescription: faker.helpers.randomize([faker.random.word(), undefined]), valueLabel: faker.random.word(), valueDescription: faker.helpers.randomize([faker.random.word(), undefined]), value: faker.random.word()})), status: faker.helpers.randomize(Object.values(RequestAccessGroupTargetStatus)), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}})), approvalMethod: faker.helpers.randomize([faker.helpers.randomize(Object.values(RequestAccessGroupApprovalMethod)), undefined]), accessRule: {timeConstraints: {maxDurationSeconds: faker.datatype.number()}}, requestStatus: faker.helpers.randomize(Object.values(RequestStatus)), requestReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), groupReviewers: faker.helpers.randomize([[...Array(faker.datatype.number({min: 1, max: 10}))].map(() => (faker.random.word())), undefined]), finalTiming: faker.helpers.randomize([{startTime: faker.random.word(), endTime: faker.random.word()}, undefined])})), requestedBy: {id: faker.random.word(), firstName: faker.random.word(), lastName: faker.random.word(), email: faker.random.word(), picture: faker.helpers.randomize([faker.random.word(), undefined])}, requestedAt: faker.random.word(), status: faker.helpers.randomize(Object.values(RequestStatus)), targetCount: faker.datatype.number()})

export const getGetGroupTargetInstructionsMock = () => ({instructions: {instructions: faker.helpers.randomize([faker.random.word(), undefined])}})

export const getDefaultMSW = () => [
rest.delete('*/api/v1/admin/access-rules/:ruleId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.post('*/api/v1/admin/target-groups/:id/resources/:resourceType/filters', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminFilterTargetGroupResourcesMock()),
        )
      }),rest.get('*/api/v1/admin/target-groups/:id/routes', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getAdminListTargetRoutesMock()),
        )
      }),rest.get('*/api/v1/entitlements', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserListEntitlementsMock()),
        )
      }),rest.get('*/api/v1/entitlements/targets', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserListEntitlementTargetsMock()),
        )
      }),rest.get('*/api/v1/preflight/:preflightId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserGetPreflightMock()),
        )
      }),rest.post('*/api/v1/preflight', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserRequestPreflightMock()),
        )
      }),rest.get('*/api/v1/reviews', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserListReviewsMock()),
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
ctx.json(getUserPostRequestsMock()),
        )
      }),rest.get('*/api/v1/requests/:requestId', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getUserGetRequestMock()),
        )
      }),rest.get('*/api/v1/targets/:targetId/access-instructions', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetGroupTargetInstructionsMock()),
        )
      }),]
