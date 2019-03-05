mocha.setup('bdd')

addEventListener('load', (event) => {
  mocha.checkLeaks()
  mocha.run()
})
