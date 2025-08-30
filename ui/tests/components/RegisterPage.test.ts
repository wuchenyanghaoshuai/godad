import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import RegisterPage from '@/pages/RegisterPage.vue'
import { useAuthStore } from '@/stores/auth'

// Mock路由
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/register', component: RegisterPage },
    { path: '/login', component: { template: '<div>Login</div>' } }
  ]
})

describe('RegisterPage', () => {
  let wrapper: any
  let authStore: any

  beforeEach(() => {
    // 设置Pinia
    setActivePinia(createPinia())
    authStore = useAuthStore()
    
    // Mock authStore.register方法
    vi.spyOn(authStore, 'register').mockResolvedValue({})
    
    wrapper = mount(RegisterPage, {
      global: {
        plugins: [router],
        stubs: {
          RouterLink: true
        }
      }
    })
  })

  describe('密码强度验证', () => {
    it('应该显示密码强度要求', async () => {
      const passwordInput = wrapper.find('input[type="password"]')
      await passwordInput.setValue('test')
      
      // 验证实时提示是否显示
      expect(wrapper.text()).toContain('至少8位字符')
      expect(wrapper.text()).toContain('包含至少一个小写字母')
      expect(wrapper.text()).toContain('包含至少一个大写字母')
    })

    it('密码长度不足时应该显示错误', async () => {
      const passwordInput = wrapper.find('input[type="password"]')
      await passwordInput.setValue('123')
      
      // 查找长度要求指示器
      const lengthIndicator = wrapper.find('[data-testid="length-indicator"]') || wrapper.text()
      expect(lengthIndicator).toContain('○') // 未满足的指示器
    })

    it('符合要求的密码应该显示成功状态', async () => {
      const passwordInput = wrapper.find('input[type="password"]')
      await passwordInput.setValue('TestPassword123')
      
      // 检查是否所有要求都满足
      const indicators = wrapper.findAll('.text-green-600')
      expect(indicators.length).toBeGreaterThan(0)
    })
  })

  describe('注册表单提交', () => {
    it('应该验证必填字段', async () => {
      const submitButton = wrapper.find('button[type="submit"]')
      
      // 尝试提交空表单
      await submitButton.trigger('click')
      
      // 表单应该阻止提交
      expect(authStore.register).not.toHaveBeenCalled()
    })

    it('密码不一致时应该显示错误', async () => {
      // 填写表单数据
      await wrapper.find('input[id="username"]').setValue('testuser')
      await wrapper.find('input[id="email"]').setValue('test@example.com')
      await wrapper.find('input[id="nickname"]').setValue('测试用户')
      await wrapper.find('input[id="password"]').setValue('TestPassword123')
      await wrapper.find('input[id="confirmPassword"]').setValue('DifferentPassword123')
      await wrapper.find('input[type="checkbox"]').setChecked()
      
      const submitButton = wrapper.find('button[type="submit"]')
      await submitButton.trigger('click')
      
      // 应该显示密码不一致的错误
      await wrapper.vm.$nextTick()
      expect(wrapper.text()).toContain('两次输入的密码不一致')
    })

    it('有效数据应该成功注册', async () => {
      // 填写完整的有效表单数据
      await wrapper.find('input[id="username"]').setValue('testuser')
      await wrapper.find('input[id="email"]').setValue('test@example.com')
      await wrapper.find('input[id="nickname"]').setValue('测试用户')
      await wrapper.find('input[id="password"]').setValue('TestPassword123')
      await wrapper.find('input[id="confirmPassword"]').setValue('TestPassword123')
      await wrapper.find('input[type="checkbox"]').setChecked()
      
      const submitButton = wrapper.find('button[type="submit"]')
      await submitButton.trigger('click')
      
      await wrapper.vm.$nextTick()
      
      // 验证注册方法被调用
      expect(authStore.register).toHaveBeenCalledWith({
        username: 'testuser',
        email: 'test@example.com',
        nickname: '测试用户',
        password: 'TestPassword123'
      })
    })
  })

  describe('用户体验', () => {
    it('应该显示随机昵称生成按钮', () => {
      const generateButton = wrapper.find('[data-testid="generate-nickname"]') || 
                           wrapper.find('button:contains("随机生成")')
      expect(generateButton.exists()).toBe(true)
    })

    it('点击随机生成应该更新昵称', async () => {
      const generateButton = wrapper.find('[data-testid="generate-nickname"]') ||
                           wrapper.find('button')
      const nicknameInput = wrapper.find('input[id="nickname"]')
      
      const originalValue = nicknameInput.element.value
      await generateButton.trigger('click')
      
      // 昵称应该被更新
      await wrapper.vm.$nextTick()
      const newValue = nicknameInput.element.value
      expect(newValue).not.toBe(originalValue)
      expect(newValue.length).toBeGreaterThan(0)
    })
  })
})