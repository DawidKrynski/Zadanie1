// Support file for Cypress e2e tests

Cypress.on('uncaught:exception', () => {
  return false
})

const blockBacktrace = (win) => {
  const originalFetch = win.fetch.bind(win)

  win.fetch = (input, init) => {
    const url = typeof input === 'string' ? input : input?.url

    if (url?.includes('backtrace.io')) {
      return Promise.resolve(new win.Response('', { status: 204 }))
    }

    return originalFetch(input, init)
  }
}

Cypress.Commands.overwrite('visit', (originalFn, url, options = {}) => {
  const originalOnBeforeLoad = options.onBeforeLoad

  return originalFn(url, {
    ...options,
    onBeforeLoad(win) {
      blockBacktrace(win)
      originalOnBeforeLoad?.(win)
    },
  })
})

beforeEach(() => {
  cy.intercept('**backtrace.io/**', {
    statusCode: 204,
    body: '',
    log: false,
  })
})

// Komenda logowania przez stan aplikacji, zeby testy nie powtarzaly pelnego UI login flow.
Cypress.Commands.add('login', () => {
  cy.setCookie('session-username', 'standard_user')
  cy.visit('/', {
    onBeforeLoad(win) {
      win.document.cookie = 'session-username=standard_user; path=/'
    },
  })
  cy.window().then((win) => {
    win.history.pushState({}, '', '/inventory.html')
    win.dispatchEvent(new win.PopStateEvent('popstate'))
  })
  cy.url().should('include', '/inventory.html')
  cy.get('.inventory_list').should('be.visible')
})
