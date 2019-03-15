mocha.setup('bdd')
addEventListener('load', () => { mocha.checkLeaks().run() })

import './app.test.js'
import './search.test.js'
