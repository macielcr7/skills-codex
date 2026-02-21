import { describe, expect, it } from 'vitest'
import { z } from 'zod'

import { EntityValidationError } from '#/domain/shared/error/entity_validation_error'

describe('EntityValidationError', () => {
  it('sets name and issues', () => {
    const schema = z.object({ title: z.string().min(1) })
    const result = schema.safeParse({ title: '' })
    if (result.success) {
      throw new Error('expected schema validation to fail')
    }

    const err = new EntityValidationError('Video', result.error.issues)
    expect(err.name).toBe('EntityValidationError')
    expect(err.message).toBe('Video validation failed')
    expect(err.issues.length).toBeGreaterThan(0)
  })
})

