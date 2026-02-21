import { describe, expect, it } from 'vitest'

import { HealthController } from '#/infrastructure/nest/common/controllers/health.controller'

describe('HealthController', () => {
  it('returns ok status', () => {
    const controller = new HealthController()
    expect(controller.get()).toEqual({ status: 'ok' })
  })
})
