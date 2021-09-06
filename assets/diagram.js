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
      select: (t) => this._select(t),
      setState: (s) => this._setState(s),
    }
  }

  _setState(s) {
    console.log("Set state", s, " for ", this.selected)
    if (this.selected == undefined) {
      return null
    }
    this.selected.state = s.Id
    eve("state", this, {t: this.selected})
    return this.selected
  }

  _createTeeth(range, rowidx, isKid) {
    let ts = new Array()
    range.forEach((v, colidx,_) => 
      ts.push(new Tooth(v, isKid ? colidx + 3 : colidx, rowidx, isKid))
    )
    return ts
  }


  _select(t) {
    console.log("in select ", t)
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
  constructor(name, wi, hi, api) {
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

    this.api = api
    this.teethToGraphics = new Map()
    eve.on("state", 
      (o, e) => {
        console.log("eve handler: ", o, e)
        this.redraw(o.t)
      })
  }

  create(teeth) {
    for (let t of teeth.all) {
      this.teethToGraphics.set(t, this.create_graphics(t))
      this.redraw(t)
    }
  }

  redraw(t) {
    console.log(`redraw ${t}`)
    let selcols = t.isSelected ? this.colors.sel : this.colors.norm
    let width = t.isSelected ? 4 : 1
    let cols = t.isKid ? selcols.kid : selcols.adult
    let toothrect = this.teethToGraphics.get(t)
    toothrect.attr("fill",  t.state === 0 ? cols.fill : this.api.getColor(t.state)) 
    toothrect.attr("stroke", cols.stroke)
    toothrect.attr("stroke-width", width)

  }

  create_graphics(t) {
    let x = this.space + t.col * (this.space + this.cell)
    let y = this.space + t.row * (this.space + this.cell)
    var toothrect = this.r.rect(x, y, this.cell, this.cell, 5)
    
    let toothNum = this.r.text(x + this.cell/2, y + this.cell/ 2, t.num)
    toothNum.attr("fill", "#222")


    toothrect.click(() => {
      console.log("in click for tooth ", t)
      let forUpdate = this.api.select(t)
      forUpdate.forEach(t => this.redraw(t))
    })

    console.log(`created graphics ${t} at ${x}, ${y}`)
    return toothrect
  }

}

class States{
  constructor(s) {
    this.states = s
    this.colors = new Map(
      s.map(
        (el, idx) => 
          [el.Id, Raphael.hsl2rgb(idx * (360/this.states.length), 50, 60)]))
    this.colors.set(0, "white")
    console.log("States created with s = ", s , " and colors ", this.colors)

    this.api = {
      getColor: (s) => this._getColors(s)
    }
  }

  _getColors(s) {
    return this.colors.get(s)
  }
}

function *pos_iter(cols, width, dy) {
  var curr_col = 0
  var curr_y = 0
  var curr_x = 0
  let dx = width/cols
  while (true) {
    yield [curr_x, curr_y]
    curr_col  = (curr_col + 1) % cols
    if (curr_col === 0) {
      curr_y += dy
    }
    curr_x = curr_col * dx
  }
}

function get_n(iter, n) {
  let res = []
  for (let i = 0; i < n; i++) {
    let { value, done } = iter.next()
    res.push(done ? null : value)
  }
  return res
}

class StatesGrid{
  constructor(name, wi, hi, api){
    console.log(arguments)
    this.r = Raphael(document.getElementById(name), wi, hi)
    this.api = api
    this.stateToGraphics = new Map()
    this.height = hi
    this.width = wi
  }


  create(stateModel) {
    let cell_x_count = 3
    let cell_width = this.width / cell_x_count
    let cell_height = 50
    let rect_width = 30
    let rect_height = 30
    
  
    let states = stateModel.states
    let positions = get_n(pos_iter(cell_x_count, this.width, cell_height), states.length)
    
    for (let i = 0; i < states.length; i++) {
      let [x, y] = positions[i]
      let s = states[i]
      let stateGr = this.r.set()

      let bg =this.r.rect(x, y, cell_width, cell_height)
      bg.attr("fill", "white")
      bg.attr("stroke", "none")
      let rect = this.r.rect(
        x + 0.5 * rect_width, 
        y + (cell_height - rect_height)/2, 
        rect_width, 
        rect_height)
      let t =  this.r.text(
        x + 1.8 * rect_width, 
        y + (cell_height)/2, 
        //y + rect_height, 
        s.Name)
      t.attr("text-anchor", "start")
      t.attr("font-size", 12)
      rect.attr("fill", this.api.getColor(s.Id))
      rect.attr("stroke", "#bbb")

      stateGr.push(bg)
      stateGr.push(rect)
      stateGr.push(t)
      stateGr.click(() =>{
        console.log(`clicked stateGr for id=${s.Id}`)
        this.api.setState(s)
      })
      stateGr.forEach(el =>
        el.hover(()=>{
          bg.attr("fill", "#ddd")
        }, ()=> {
          bg.attr("fill", "white")
        }, null, null) 
      )
      y += cell_height
      this.stateToGraphics.set(s, stateGr)
    }
  }
}
export class App {
  constructor(toothgrid_data, states_data) {
    this.teeth = new Teeth()
    let {states} = states_data
    this.stateModel = new States(states)
    this.api = {...this.stateModel.api, ...this.teeth.api}
    console.log("App.api is: ", this.api)
    this.init_toothgrid(toothgrid_data)
    this.init_statesgrid(states_data)

  }

  init_toothgrid(toothgrid_data) {
    let { name, wi, hi } = toothgrid_data
    this.grid = new TeethGrid(name, wi, hi, this.api)
    this.grid.create(this.teeth)
  }
  init_statesgrid(states_data) {
    let {name, wi, hi, states} = states_data
    this.stateModel = new States(states)
    this.statesGrid = new StatesGrid(name, wi, hi, this.api)

    this.statesGrid.create(this.stateModel)
  }
}

