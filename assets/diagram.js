class Tooth {

  constructor(num, col, row, iskid) {
    this.num = num
    this.col = col
    this.row = row
    this.isKid = iskid
    this.state = 0
    this.isSelected = false
    this.init_sides()
  }
  
  init_sides() {
    this.sides = [0, 0, 0, 0, 0]
  }

  toString() {
    return `[${this.num}] (${this.isSelected ? "selected" : ""})[${this.col}x${this.row}] state: ${this.state} sides: ${this.sides}`
  }


}

class Changer {
  MAX_SIDES = 5
  constructor(states) {
    this.states = states
  }

  check_state(t, s) {
    let max = this.states.length
    if ((s < 0) || (s > max) ){
      throw Error(`State ${s} for tooth ${t} invalid. Should be in range [0, ${max}]`)
    }
    return s
  }
  // TODO
  check_side(t, s) {
    let max = Changer.MAX_SIDES
    if ((s < 0) || (s > max)) {
      throw Error(`Side ${s} for tooth ${t} invalid. Should be in range [0, ${max}]`)
    }
    return s
  }
  // TODO better value filtering if new values are added but old "removed"
  check_value(t, s) {
    let v = this.states.map(s => s.Id)
    if (!v.includes(s)) {
      throw Error(`Value ${s} for tooth ${t} invalid. Should be in  ${v}`)
    }
    return s
  }


  change(t, {StateId, IsWhole,  ToothSide}) {
    t.state = this.check_state(t, StateId)
    if (!IsWhole) {
      t.sides[this.check_side(t, ToothSide)] = this.check_value(StateId)
    } else {
      t.init_sides()
    }
  }
}

function range(a, b) {
  let res = new Array()
  let fn =  (a <= b) ? Array.prototype.push : Array.prototype.unshift
  for (let i = Math.min(a, b); i <= Math.max(a, b); i++) {
    fn.apply(res, [i]) 
  }
  return res
  

}
class Teeth {
  constructor() {
    let adult_up =  this._createTeeth(range(18, 11).concat(range(21, 28)), 0, false)
    let adult_down =  this._createTeeth(range(48, 41).concat(range(31, 38)), 3, false)
    let kid_up =  this._createTeeth(range(55, 51).concat(range(61, 65)), 1, true)
    let kid_down = this._createTeeth(range(85, 81).concat(range(71, 75)), 2, true)
    
    this.all = adult_up.concat(adult_down).concat(kid_up).concat(kid_down)
    this.selected = undefined

    this.api = {
      select: (t) => this._select(t)
    }
  }

  _createTeeth(range, rowidx, isKid) {
    let ts = new Array()
    range.forEach((v, colidx,_) => 
      ts.push(new Tooth(v, isKid ? colidx + 3 : colidx, rowidx, isKid))
    )
    return ts
  }

  

  _select(t) {
    let forUpdate = []
    if (this.selected != undefined) {
      forUpdate.push(this.selected)
      this.selected.isSelected = false
    }
    this.selected = t
    t.isSelected = true
    forUpdate.push(t)
    return forUpdate
  }


}

class TeethGrid {
  constructor(name, wi, hi, teeth_api) {
    this.r = Raphael(document.getElementById(name), wi, hi)
    this.space = 5
    this.cell = (Math.max(wi, hi) / 16) - 2 * this.space
    this.colors =  {
      norm : {
        kid: {fill: "#fffbeb", stroke: "#48c774"},
        adult:{fill: "#edf6fc", stroke: "#3298dc"} 
      },

      sel : {
        kid: {fill: "#feecf0", stroke: "#f14668"},
        adult:{fill: "#feecf0", stroke: "#f14668"} 
      }
    }

    this.api = teeth_api
    this.teethToGraphics = new Map()
  }

  create(teeth) {
    for (let t of teeth.all) {
      this.teethToGraphics.set(t, this.create_graphics(t))
      this.redraw(t)
    }
  }

  redraw(t) {
    console.log(`redraw ${t}`)
    let cols = t.isSelected ? this.colors.sel : this.colors.norm
    let toothrect = this.teethToGraphics.get(t)
    toothrect.attr("fill", t.isKid ? cols.kid.fill : cols.adult.fill)
    toothrect.attr("stroke", t.isKid ? cols.kid.stroke : cols.adult.stroke)
  }

  create_graphics(t) {
    let x = this.space + t.col * (this.space + this.cell)
    let y = this.space + t.row * (this.space + this.cell)
    var toothrect = this.r.rect(x, y, this.cell, this.cell)
    
    let toothNum = this.r.text(x + this.cell/2, y + this.cell/ 2, t.num)
    toothNum.attr("fill", "#222")


    toothrect.click(() => {
      let forUpdate = this.api.select(t)
      forUpdate.forEach(t => this.redraw(t))
    })

    console.log(`created graphics ${t} at ${x}, ${y}`)
    return toothrect
  }

}

export class App {
  constructor(name, wi, hi, states) {
    this.teeth = new Teeth()
    this.grid = new TeethGrid(name, wi, hi, this.teeth.api)
    this.grid.create(this.teeth)
    this.changer = new Changer(states)
  }
  
}
