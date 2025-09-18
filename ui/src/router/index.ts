import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import LoginPage from '@/pages/LoginPage.vue'
import RegisterPage from '@/pages/RegisterPage.vue'
import ForgotPasswordPage from '@/pages/ForgotPasswordPage.vue'
import ResetPasswordPage from '@/pages/ResetPasswordPage.vue'
import UserCenter from '@/pages/UserCenter.vue'
import UserProfilePage from '@/pages/UserProfilePage.vue'
import CategoryManagePage from '@/pages/CategoryManagePage.vue'
import ArticleListPage from '@/pages/ArticleListPage.vue'
import ArticleDetailPage from '@/pages/ArticleDetailPage.vue'
import ArticleCreatePage from '@/pages/ArticleCreatePage.vue'
import ArticleEditPage from '@/pages/ArticleEditPage.vue'
import SearchPage from '@/pages/SearchPage.vue'
import NotFoundPage from '@/pages/NotFoundPage.vue'
import AboutView from '@/views/AboutView.vue'
import AdminDashboard from '@/pages/AdminDashboard.vue'
import AdminLoginPage from '@/pages/AdminLoginPage.vue'
import NotificationsPage from '@/pages/NotificationsPage.vue'
import AboutPage from '@/pages/AboutPage.vue'
import ContactPage from '@/pages/ContactPage.vue'
import PrivacyPage from '@/pages/PrivacyPage.vue'
import TermsPage from '@/pages/TermsPage.vue'
import HelpPage from '@/pages/HelpPage.vue'
import CommunityPage from '@/pages/CommunityPage.vue'
import ResourcesPage from '@/pages/ResourcesPage.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    // 如果有保存的位置（比如浏览器前进后退），使用保存的位置
    if (savedPosition) {
      return savedPosition
    }
    // 否则滚动到页面顶部
    return { top: 0 }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomePage,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterPage,
      meta: { requiresGuest: true }
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: ForgotPasswordPage,
      meta: { requiresGuest: true }
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: ResetPasswordPage,
      meta: { requiresGuest: true }
    },
    {
      path: '/user-center',
      name: 'user-center',
      component: UserCenter,
      meta: { requiresAuth: true }
    },
    {
      path: '/users/:username',
      name: 'UserProfile',
      component: UserProfilePage,
      meta: { requiresFromArticle: true }
    },
    {
      path: '/notifications',
      name: 'Notifications',
      component: NotificationsPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/about',
      name: 'about',
      component: AboutPage,
    },
    {
      path: '/contact',
      name: 'contact',
      component: ContactPage,
    },
    {
      path: '/privacy',
      name: 'privacy',
      component: PrivacyPage,
    },
    {
      path: '/terms',
      name: 'terms',
      component: TermsPage,
    },
    {
      path: '/help',
      name: 'help',
      component: HelpPage,
    },
    {
      path: '/admin/login',
      name: 'AdminLogin',
      component: AdminLoginPage,
      meta: { requiresGuest: true }
    },
    {
      path: '/management-dashboard',
      name: 'AdminDashboard',
      component: AdminDashboard,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/management-dashboard/categories',
      name: 'CategoryManage',
      component: CategoryManagePage,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/articles',
      name: 'ArticleList',
      component: ArticleListPage
    },
    {
      path: '/articles/create',
      name: 'ArticleCreate',
      component: ArticleCreatePage,
      meta: { requiresAuth: true }
    },
    {
      path: '/articles/:id',
      name: 'ArticleDetail',
      component: ArticleDetailPage
    },
    {
      path: '/articles/:id/edit',
      name: 'ArticleEdit',
      component: ArticleEditPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/search',
      name: 'Search',
      component: SearchPage
    },
    {
      path: '/categories',
      name: 'Categories',
      component: ArticleListPage // 使用文章列表页面，通过查询参数区分
    },
    {
      path: '/community',
      name: 'Community',
      component: CommunityPage
    },
    {
      path: '/resources',
      name: 'Resources',
      component: ResourcesPage
    },
    {
      path: '/404',
      name: 'NotFound',
      component: NotFoundPage
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/404'
    }
  ],
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  try {
    // 如果还没有初始化认证状态，先初始化
    if (!authStore.token && !authStore.user) {
      await authStore.initAuth()
    }
    
    // 检查用户认证状态
    if (authStore.token && !authStore.isAuthenticated) {
      await authStore.checkAuth()
    }
    
    // 需要认证的路由
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      // 管理员路由跳转到管理员登录页面
      if (to.meta.requiresAdmin) {
        next({
          path: '/admin/login',
          query: { redirect: to.fullPath }
        })
      } else {
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
      }
      return
    }
    
    // 需要管理员权限的路由
    if (to.meta.requiresAdmin && (!authStore.isAuthenticated || !authStore.isAdmin)) {
      next('/404')
      return
    }
    
    // 游客专用路由（已登录用户不能访问）
    if (to.meta.requiresGuest && authStore.isAuthenticated) {
      const redirect = to.query.redirect as string
      next(redirect || '/')
      return
    }
    
    // 验证文章ID参数
    if (to.name === 'ArticleDetail' || to.name === 'ArticleEdit') {
      const id = to.params.id
      if (!id || isNaN(Number(id))) {
        next('/404')
        return
      }
    }
    
    // 允许用户主页访问（移除限制）
    // 注释掉原有的限制逻辑，用户应该可以直接访问任何用户的主页
    /*
    if (to.meta.requiresFromArticle && to.name === 'UserProfile') {
      // 检查是否从文章相关页面或消息页面跳转，或者是页面刷新
      const isFromAllowedPage = from.name === 'ArticleDetail' ||
                               from.name === 'ArticleList' ||
                               from.name === 'home' ||
                               from.name === 'Messages' ||
                               from.path.startsWith('/articles') ||
                               from.name === null // 页面刷新时from.name为null

      if (!isFromAllowedPage) {
        // 如果不是从允许的页面跳转，重定向到首页
        next('/')
        return
      }
    }
    */
    
    // 验证用户名参数
    if (to.name === 'UserProfile') {
      const username = to.params.username
      if (!username || typeof username !== 'string' || username.trim() === '') {
        next('/404')
        return
      }
    }
    
    next()
  } catch (error) {
    console.error('路由守卫错误:', error)
    // 发生错误时重定向到首页
    if (to.path !== '/') {
      next('/')
    } else {
      next()
    }
  }
})

// 路由错误处理
router.onError((error) => {
  console.error('路由错误:', error)
  // 可以在这里添加错误上报逻辑
})

export default router
