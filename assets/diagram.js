class Web {
}
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
    eve("select", this, forUpdate)
  }


}

class TeethGrid {
  constructor(name, wi, hi, api) {
    this.r = Raphael(document.getElementById(name), wi, hi)
    this.space = 5
    this.cell = (Math.max(wi, hi) / 16) - 2 * this.space
    this.colors =  {
      norm : {
        kid: {fill: "#fffbeb", stroke: "#48c774", strokewidth: 1},
        adult:{fill: "#edf6fc", stroke: "#3298dc", strokewidth: 1}
      },

      sel : {
        kid: {fill: "#feecf0", stroke: "#f14668", strokewidth: 4},
        adult:{fill: "#feecf0", stroke: "#f14668", strokewidth: 4}
      }
    }

    this.api = api
    this.teethToGraphics = new Map()

    eve.on("state",
      (o) => {
        console.log("eve handler state: ", o)
        this.redraw(o.t)
      })

    eve.on("select",
      (o) => {
        console.log("eve handler select", o)
        o.forEach(t => this.redraw(t))
      })
  }

  create(teeth) {
    for (let t of teeth.all) {
      let gr = this.create_graphics(t)
      this.teethToGraphics.set(t, gr)
      this.redraw(t)
      gr.hover(()=>{
        gr.attr("stroke", "#f14668")
      }, ()=> {
        let {stroke} = this.getFillStroke(t)
        gr.attr("stroke", stroke)
      }, null, null)
    }
  }

  getFillStroke(t) {
    let selcols = t.isSelected ? this.colors.sel : this.colors.norm
    let cols = t.isKid ? selcols.kid : selcols.adult
    let fill = t.state === 0 ? cols.fill : this.api.getColor(t.state)
    return {fill, stroke: cols.stroke, strokewidth: cols.strokewidth}
  }


  setAttrs(gr, {fill, stroke, strokewidth}) {
    gr.attr("fill", fill)
    gr.attr("stroke", stroke)
    gr.attr("stroke-width", strokewidth)
  }

  redraw(t) {
    console.log(`redraw ${t}`)
    this.setAttrs(this.teethToGraphics.get(t), this.getFillStroke(t))

  }

  create_graphics(t) {
    let x = this.space + t.col * (this.space + this.cell)
    let y = this.space + t.row * (this.space + this.cell)
    var toothrect = this.r.rect(x, y, this.cell, this.cell, 10)

    let toothNum = this.r.text(x + this.cell/2, y + this.cell/ 2, t.num)
    toothNum.attr("fill", "#222")


    toothrect.click(() => {
      console.log("in click for tooth ", t)
      this.api.select(t)
    })

    console.log(`created graphics ${t} at ${x}, ${y}`)
    return toothrect
  }

}

class States{
  constructor(s) {
    this.states = s
    this.idToSttae = new Map( s.map(el => [el.Id, el]))
    console.log("States created with s = ", s , " and colors ", this.colors)

    this.api = {
      getColor: (s) => this._getColors(s),
      getPath: (s) => this._getPath(s)
    }
  }
  _getPath(s) {
    if (!this.idToState.has(s)) {
      // rect by default
      return "M0,0h100v100h-100Z"
    }
    return this.idToState.get(s).Svgpath

  }
  _getColors(s) {
    if (!this.idToState.has(s)) {
      // pink color by default
      return "pink"
    }
    return this.idToState.get(s).Color
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

      this.r.translate(
        x + 0.5 * rect_width,
        y + (cell_height - rect_height)/2)
      this.r.scale(rect_width/100, rect_height/100)
      let rect = this.r.path(this.api.getPath(s.Id))
      
      let t =  this.r.text(
        x + 1.8 * rect_width,
        y + (cell_height)/2,
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

