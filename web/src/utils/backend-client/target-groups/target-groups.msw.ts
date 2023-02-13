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

export const getListTargetGroupDeploymentsMock = () => ({res: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), functionArn: faker.random.word(), awsAccount: faker.random.word(), awsRegion: faker.random.word(), healthy: faker.datatype.boolean(), diagnostics: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.random.word(), code: faker.random.word(), message: faker.random.word()})), activeConfig: faker.helpers.arrayElement([{
        'cle2a7baa000n52on1vfpgkun': {type: faker.random.word(), value: {}}
      }, undefined])})), next: faker.random.word()})

export const getCreateTargetGroupDeploymentMock = () => ({id: faker.random.word(), functionArn: faker.random.word(), awsAccount: faker.random.word(), awsRegion: faker.random.word(), healthy: faker.datatype.boolean(), diagnostics: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({level: faker.random.word(), code: faker.random.word(), message: faker.random.word()})), activeConfig: faker.helpers.arrayElement([{
        'cle2a7baa000o52ondxbmf6qi': {type: faker.random.word(), value: {}}
      }, undefined])})

export const getGetTargetGroupMock = () => ({id: faker.random.word(), targetSchema: {From: faker.random.word(), Schema: {
        'cle2a7bab000q52on7usgfjog': {id: faker.random.word(), title: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined]), ruleFormElement: faker.helpers.arrayElement(['INPUT','MULTISELECT','SELECT']), requestFormElement: faker.helpers.arrayElement(['SELECT']), groups: {
        'cle2a7bab000p52on5t4a37uj': {id: faker.random.word(), title: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined])}
      }}
      }}, icon: faker.random.word(), createdAt: faker.helpers.arrayElement([faker.random.word(), undefined]), updatedAt: faker.helpers.arrayElement([faker.random.word(), undefined])})

export const getListTargetGroupsMock = () => ({targetGroups: Array.from({ length: faker.datatype.number({ min: 1, max: 10 }) }, (_, i) => i + 1).map(() => ({id: faker.random.word(), targetSchema: {From: faker.random.word(), Schema: {
        'cle2a7bad000s52onhbj4bnud': {id: faker.random.word(), title: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined]), ruleFormElement: faker.helpers.arrayElement(['INPUT','MULTISELECT','SELECT']), requestFormElement: faker.helpers.arrayElement(['SELECT']), groups: {
        'cle2a7bad000r52on67t52ulp': {id: faker.random.word(), title: faker.random.word(), description: faker.helpers.arrayElement([faker.random.word(), undefined])}
      }}
      }}, icon: faker.random.word(), createdAt: faker.helpers.arrayElement([faker.random.word(), undefined]), updatedAt: faker.helpers.arrayElement([faker.random.word(), undefined])})), next: faker.helpers.arrayElement([faker.random.word(), undefined])})

export const getTargetGroupsMSW = () => [
rest.get('*/api/v1/target-group-deployments/:id', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.get('*/api/v1/target-group-deployments', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getListTargetGroupDeploymentsMock()),
        )
      }),rest.post('*/api/v1/target-group-deployments', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getCreateTargetGroupDeploymentMock()),
        )
      }),rest.get('*/api/v1/target-groups/:id', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getGetTargetGroupMock()),
        )
      }),rest.get('*/api/v1/target-groups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
ctx.json(getListTargetGroupsMock()),
        )
      }),rest.post('*/api/v1/target-groups', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),rest.post('*/api/v1/target-groups/:id/link', (_req, res, ctx) => {
        return res(
          ctx.delay(1000),
          ctx.status(200, 'Mocked status'),
        )
      }),]
