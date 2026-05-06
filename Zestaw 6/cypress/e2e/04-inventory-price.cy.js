describe('Produkty - sortowanie cen', () => {
  beforeEach(() => {
    cy.login()
  })

  it('powinien sortowac produkty po cenie rosnaco', () => {
    cy.get('[data-test="product-sort-container"]').select('lohi')
    cy.get('.inventory_item_price').first().invoke('text').then((text) => {
      const price = parseFloat(text.replace('$', ''))
      expect(price).to.equal(7.99)
    })
  })

  it('powinien sortowac produkty po cenie malejaco', () => {
    cy.get('[data-test="product-sort-container"]').select('hilo')
    cy.get('.inventory_item_price').first().invoke('text').then((text) => {
      const price = parseFloat(text.replace('$', ''))
      expect(price).to.equal(49.99)
    })
  })
})
