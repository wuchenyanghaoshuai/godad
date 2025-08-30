import { test, expect } from '@playwright/test'

test.describe('管理员功能流程', () => {
  test.beforeEach(async ({ page }) => {
    // 导航到首页
    await page.goto('http://127.0.0.1:3333')
  })

  test('管理员登录并访问后台', async ({ page }) => {
    // 步骤1: 导航到登录页面
    await page.click('text=登录')
    await expect(page).toHaveURL(/.*login/)

    // 步骤2: 输入管理员凭据
    await page.fill('input[id="username"]', 'chenyangwu')
    await page.fill('input[id="password"]', 'chenyangwu的密码') // 实际测试中需要使用真实密码
    await page.click('button[type="submit"]')

    // 步骤3: 验证登录成功，跳转到首页
    await expect(page).toHaveURL('http://127.0.0.1:3333/')
    
    // 步骤4: 点击用户头像打开菜单
    await page.click('[data-testid="user-menu-button"]') // 需要在实际组件中添加data-testid
    
    // 步骤5: 验证管理员菜单项存在
    await expect(page.locator('text=后台管理')).toBeVisible()
    
    // 步骤6: 点击后台管理
    await page.click('text=后台管理')
    
    // 步骤7: 验证跳转到管理员后台
    await expect(page).toHaveURL(/.*management-dashboard/)
    
    // 步骤8: 验证统计数据显示
    await expect(page.locator('[data-testid="article-count"]')).toBeVisible()
    await expect(page.locator('[data-testid="user-count"]')).toBeVisible()
    await expect(page.locator('[data-testid="category-count"]')).toBeVisible()
    await expect(page.locator('[data-testid="comment-count"]')).toBeVisible()
    
    // 步骤9: 验证数据不为空
    const articleCount = await page.locator('[data-testid="article-count"]').textContent()
    expect(parseInt(articleCount || '0')).toBeGreaterThan(0)
  })

  test('普通用户无法看到管理员菜单', async ({ page }) => {
    // 创建一个普通用户账号进行测试
    // 这里假设已有测试用户账号
    
    // 步骤1: 登录普通用户
    await page.click('text=登录')
    await page.fill('input[id="username"]', 'testuser')
    await page.fill('input[id="password"]', 'TestPassword123')
    await page.click('button[type="submit"]')
    
    // 步骤2: 验证登录成功
    await expect(page).toHaveURL('http://127.0.0.1:3333/')
    
    // 步骤3: 点击用户头像
    await page.click('[data-testid="user-menu-button"]')
    
    // 步骤4: 验证没有管理员菜单
    await expect(page.locator('text=后台管理')).not.toBeVisible()
    
    // 步骤5: 验证只有普通用户菜单项
    await expect(page.locator('text=个人中心')).toBeVisible()
    await expect(page.locator('text=发布文章')).toBeVisible()
  })

  test('直接访问管理员URL应该被拦截', async ({ page }) => {
    // 步骤1: 未登录状态直接访问管理员页面
    await page.goto('http://127.0.0.1:3333/management-dashboard')
    
    // 步骤2: 应该被重定向到登录页面
    await expect(page).toHaveURL(/.*login/)
    
    // 步骤3: 登录普通用户
    await page.fill('input[id="username"]', 'testuser')
    await page.fill('input[id="password"]', 'TestPassword123')
    await page.click('button[type="submit"]')
    
    // 步骤4: 再次尝试访问管理员页面
    await page.goto('http://127.0.0.1:3333/management-dashboard')
    
    // 步骤5: 应该被重定向到404或首页
    await expect(page).toHaveURL(/.*404|.*\/$/)
  })
})

test.describe('用户注册流程', () => {
  test('完整注册流程', async ({ page }) => {
    await page.goto('http://127.0.0.1:3333')
    
    // 步骤1: 点击注册
    await page.click('text=注册')
    await expect(page).toHaveURL(/.*register/)
    
    // 步骤2: 填写注册表单
    const timestamp = Date.now()
    await page.fill('input[id="username"]', `testuser${timestamp}`)
    await page.fill('input[id="email"]', `test${timestamp}@example.com`)
    await page.fill('input[id="nickname"]', `测试用户${timestamp}`)
    
    // 步骤3: 测试密码强度验证
    await page.fill('input[id="password"]', '123')
    await expect(page.locator('text=至少8位字符')).toHaveClass(/text-gray-400/)
    
    await page.fill('input[id="password"]', 'TestPassword123')
    await expect(page.locator('text=至少8位字符')).toHaveClass(/text-green-600/)
    
    // 步骤4: 确认密码
    await page.fill('input[id="confirmPassword"]', 'TestPassword123')
    
    // 步骤5: 同意条款
    await page.check('input[type="checkbox"]')
    
    // 步骤6: 提交注册
    await page.click('button[type="submit"]')
    
    // 步骤7: 验证跳转到登录页面并显示成功消息
    await expect(page).toHaveURL(/.*login/)
    await expect(page.locator('text=注册成功')).toBeVisible()
    
    // 步骤8: 使用新账号登录
    await page.fill('input[id="username"]', `testuser${timestamp}`)
    await page.fill('input[id="password"]', 'TestPassword123')
    await page.click('button[type="submit"]')
    
    // 步骤9: 验证登录成功
    await expect(page).toHaveURL('http://127.0.0.1:3333/')
    
    // 步骤10: 验证用户菜单不包含管理员选项
    await page.click('[data-testid="user-menu-button"]')
    await expect(page.locator('text=后台管理')).not.toBeVisible()
  })

  test('密码强度验证', async ({ page }) => {
    await page.goto('http://127.0.0.1:3333/register')
    
    // 测试各种密码强度场景
    const passwordTests = [
      { password: '123', shouldFail: true, reason: '长度不足' },
      { password: 'password', shouldFail: true, reason: '没有大写字母' },
      { password: 'PASSWORD', shouldFail: true, reason: '没有小写字母' },
      { password: 'Password123', shouldFail: false, reason: '符合要求' }
    ]
    
    for (const test of passwordTests) {
      await page.fill('input[id="password"]', test.password)
      
      if (test.shouldFail) {
        // 提交按钮应该禁用或显示错误
        const submitButton = page.locator('button[type="submit"]')
        const isDisabled = await submitButton.getAttribute('disabled')
        expect(isDisabled).not.toBeNull()
      }
    }
  })
})