import { describe, expect, it } from 'vitest'

import { UUIDGenerator } from '#/infrastructure/nest/common/services/uuid_generator'

describe('UUIDGenerator', () => {
  it('generates uuid v4 strings', () => {
    const generator = new UUIDGenerator()

    const id1 = generator.newID()
    const id2 = generator.newID()

    expect(id1).toMatch(/^[0-9a-f-]{36}$/)
    expect(id2).toMatch(/^[0-9a-f-]{36}$/)
    expect(id1).not.toBe(id2)
  })
})
