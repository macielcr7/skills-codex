import type { AddressInfo } from 'node:net'

import { Test } from '@nestjs/testing'
import type { INestApplication } from '@nestjs/common'
import { describe, expect, it, afterAll, beforeAll } from 'vitest'

import { AppModule } from '#/app.module'

describe('Media/Video E2E', () => {
  let app: INestApplication
  let baseUrl: string

  beforeAll(async () => {
    const moduleRef = await Test.createTestingModule({
      imports: [AppModule],
    }).compile()

    app = moduleRef.createNestApplication()
    await app.init()
    await app.listen(0)

    const address = app.getHttpServer().address()
    if (!address || typeof address === 'string') {
      throw new Error('unexpected server address')
    }

    baseUrl = `http://127.0.0.1:${(address as AddressInfo).port}`
  })

  afterAll(async () => {
    await app.close()
  })

  it('uploads and gets a video', async () => {
    const uploadRes = await fetch(`${baseUrl}/media/videos`, {
      method: 'POST',
      headers: { 'content-type': 'application/json' },
      body: JSON.stringify({ title: 'hello' }),
    })

    expect(uploadRes.status).toBe(201)
    const uploadBody = (await uploadRes.json()) as { id: string }
    expect(uploadBody.id).toMatch(/[0-9a-f-]{36}/)

    const getRes = await fetch(`${baseUrl}/media/videos/${uploadBody.id}`)
    expect(getRes.status).toBe(200)
    const getBody = (await getRes.json()) as { id: string; title: string }
    expect(getBody).toEqual({ id: uploadBody.id, title: 'hello' })
  })

  it('returns 404 when not found', async () => {
    const res = await fetch(
      `${baseUrl}/media/videos/550e8400-e29b-41d4-a716-446655440000`,
    )
    expect(res.status).toBe(404)
  })
})
