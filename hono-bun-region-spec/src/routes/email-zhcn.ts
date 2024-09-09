import { Hono } from 'hono'

const router = new Hono()

router.get('/example', (c) => {
    console.info(c.req.method, c.req.path)
    return c.text('Hello, Hono!')
})

export default router
