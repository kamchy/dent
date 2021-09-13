import { createMyEvents  } from "./myevents.js"

function mk_range(f, t) {
  let arr = []
  let minv = Math.min(f, t)
  let maxv = Math.max(f, t)
  let fn = f < t ? Array.prototype.push  : Array.prototype.unshift
  for (let i = minv; i <= maxv; i++) {
    fn.call(arr, i)
  }
  return arr
}

let TYPE_DATA = {
  "adult": {attr: {stroke: '#333', fill: 'Azure'}, nskip: 0, nskip_right: 8},
  "kid": {attr: {stroke: '#444', fill: 'Snow'}, nskip: 3, nskip_right: 5},
}

function text_size(d, {example, fontsize}) {
  let font = {size: fontsize}
  let m = d.text(example).font(font)
  let tdx = m.bbox().width
  let tdy = m.bbox().height
  m.remove()
  return {tdx, tdy, fontsize}
}

class States {
  constructor(states) {
    this.whole = states.filter(s => s.Whole === true)
    this.part = states.filter(s => s.Whole === false)
    console.log("whole and part: ", this.whole, this.part)
  }

  getStatesIterator(iswhole) {
    function *giterator(tab) {
      let idx = 0
      let len = tab.length
      while (true) {
        yield tab[idx % len]
        idx +=1 
      }
    }

    return iswhole ? giterator(this.whole) : giterator(this.part)
  }
}
export class App {
  constructor (data) {
    this.es = createMyEvents()
    this.tooth_drawer = new ToothDrawer(data, this.es)
    this.states_drawer = new StatesDrawer(data)
    this.sender = new Sender(data.pid, data.vid, this.es)
  }
}

class StatesDrawer {
  constructor ({name, wi, hi, states}) {
    console.log(`Name: ${name}, states`, states)
    if (states === null) {
      return
    }
    this.draw = SVG().addTo(name).size(wi, hi)
    this.d = this.data(wi, hi)
    let d = this.d

    states.sort((a, b) => a.Whole? a.Id - b.Id : 1)
    console.log(states)
    // todo: spacing, layout in grid
    states.forEach( (s, idx) => {
      let g = this.draw.group()
      g.attr({stroke: "#333", fill: "white"})

      let dx = (idx % d.cols) * d.width
      let dy = Math.floor(idx / d.cols) * d.height

      let p = g.path(s.Svgpath).fill({color: s.Color})
      p.scale(0.25, 0, 0)
      let t = g.text(s.Name)
        .font({size: d.fs, fill: "#333", stroke: "none"})
        .translate(40 + d.space, 20)

      g.translate(dx, dy)

    })
  }

  data(wi, hi) {
    let cols = 3
    let width = wi/cols
    let height = 35
    let space = 2
    let {tdx, tdy, fs} = text_size(this.draw,  {example: "aaaaaaaaaaaaaa", size: 9})
    return {cols, width, height, space, tdx, tdy, fs}
  }
}

class ToothDrawer {
  constructor ({name, wi, hi, states, changes}, es) {
    this.es = es
    console.log(`Name: ${name}, states`, states, changes)
    this.draw = SVG().addTo(name).size(wi, hi)
    this.d = this.data(wi, hi, 16, 10, text_size(this.draw, {example: "12", fontsize: 12}))
    this.draw.translate(this.d.canvas_margin, 0)
    this.states = new States(states)
    this.teeth = new Map()

    let row_data = [
      {type: "adult", ranges: [[18, 11], [21, 28]]},
      {type: "kid", ranges: [[55, 51], [61, 65]]},
      {type: "kid", ranges: [[85, 81], [71, 75]]},
      {type: "adult", ranges: [[48, 41], [31, 38]]}]
    row_data.forEach((rd, idx) => this.draw_row(idx, rd))
  }

  data(wi, hi, n_adult, n_kid, textsize) {
    let canvas_margin = 10
    let rowcount = 4
    let spacex = (wi - canvas_margin*2)/n_adult
    let spacey = (hi - canvas_margin*2)/rowcount
    let margin = 3
    let rectx = spacex - margin
    let recty = spacey - margin
    return { n_adult, n_kid, spacex, spacey, margin, canvas_margin, rectx, recty, textsize}
  }

  draw_row(idx, rd) {
    let type_data = TYPE_DATA[rd.type]
    let dx = type_data.nskip * this.d.spacex
    let dy = idx * this.d.spacey
    let attr = type_data.attr
    let dx2 = (type_data.nskip_right) * this.d.spacex + this.d.margin

    this.row({range: mk_range(...rd.ranges[0]), dx, dy, attr})
    this.row({range: mk_range(...rd.ranges[1]), dx : dx + dx2, dy, attr})
  }

  row({range, dx, dy, attr }) {
    let g = this.draw.group()
    g.translate(dx, dy)
    range.map((toothnum, idx)=> {
      let t = this.tooth(g, toothnum, idx, attr)
      this.append_tooth(toothnum, t)
    })
    return g
  }

  onclick(v) {
    return (e) => this.updateState(v)  
  }

  updateState(v) {
    let toothside = v.i
    let toothnum = v.toothnum
    let switcher = this.teeth.get(toothnum)
    let stateid
    if (toothside !== undefined) {
      stateid = switcher.updatePart(toothside)
    } else {
      stateid = switcher.updateWhole()
    }

    this.es.emit("update", 
      {toothnum, stateid, toothside})

  }

  tooth(g, toothnum, idx, attr) {
    let d = this.d
    let tdx = d.textsize.tdx/2
    let tdy = d.textsize.tdy
    let tg = g.group().data({id: toothnum})
      .translate(idx*d.spacex, 0)
    let bg = tg.rect(d.rectx, d.recty).attr({fill: "Snow"})
    bg.on("mouseover",
      function(){
         this.attr({fill: "MistyRose"})
      })
    bg.on("mouseout",
      function(){
         this.attr({fill: "Snow"})
      })



      tg.text(toothnum)
        .attr({fill: attr.stroke})
        .font({size: d.textsize.fontsize})
        .move(d.rectx/2 - tdx, 0)

    let elems = tg.group()
    elems.translate(0, 1.5*tdy)

    let c = elems.circle(Math.min(d.rectx, d.recty) - tdy)
      .attr(attr)
      .move(tdy/2, 0)
      .hide()
      
    let parts = this.mk_parts(elems, d.rectx, d.recty - tdy*3)

    bg.on("click", this.onclick({toothnum, parts}))

    return new Map([["circle", c], ["parts", parts]])

  }

  mk_parts(tg, wi, hi) {
    let toothnum = tg.parent().data("id")
    let [v5, v6, v7, v8] = [[wi/3, hi/3], [2*wi/3, hi/3], [2*wi/3, 2*hi/3], [wi/3, 2*hi/3]]
    let [v1, v2, v3, v4] = [[0, 0], [wi, 0], [wi, hi], [0, hi]]
    let polygons_vertices = [[v1, v5, v8, v4], [v1, v2, v6, v5], [v2, v3, v7, v6], [v3, v4, v8, v7], [v5, v6, v7, v8]]
    let ps = tg.group()
    polygons_vertices.forEach((pv, i) => {
      ps.polygon().plot(pv)
        .on("click", this.onclick({toothnum, i}))
        .stroke({
          color: "SlateGray", 
          width: 2, 
          linecap: "round", 
          linejoin: "round"}
      ).fill({color: "white"})
    })

    return ps
  }

  append_tooth(toothnum, t) {
    this.teeth.set(toothnum, new Switcher(this.states, t))
  }



}

class Switcher{
  constructor (states, t) {
    this.t = t
    this.states = states
    this.stateIdx = 0
    this.partscount = 5
    this.partsIdx = []
    for (let i = 0; i < this.partscount; i++) {
      this.partsIdx.push(0)
    }
    console.log("Switcher states and t", states, t)
    this.c = this.t.get("circle")
    this.p = this.t.get("parts")
  }

  updateWhole() {
    this.stateIdx = (this.stateIdx + 1) % this.states.whole.length
    let s = this.states.whole[this.stateIdx]
    if (s.Id !== 0) {
    this.c = this.c.replace(this.c.parent().path(s.Svgpath).fill({color: "#777"}).stroke({color: "black"}).scale(0.5, 0, 0).move(4, 0))
    this.p.hide()
    this.c.show()
    } else {
    this.p.show()
    this.c.hide()
    }
    return s.Id
  }

  updatePart(partid) {
    let nextIdx = (this.partsIdx[partid] + 1) % this.states.part.length
    this.partsIdx[partid] = nextIdx 
    let s = this.states.part[nextIdx]
    console.log("part update: ", nextIdx, this.states.part, s)
    this.p.get(partid).attr({fill: s.Color})
    this.p.show()
    this.c.hide()
    return s.Id
  }
}

class Sender {
  constructor(pid, vid, es) {
    this.pid = pid
    this.vid = vid
    console.log(`Sender constructor: pid= ${this.pid} vid=${this.vid}`)
    this.es = es
    es.on("update", (data)=>this.sendUpdate({vid: this.vid, pid: this.pid, ...data }))
  }
  sendUpdate(data) {
    console.log("sendUpdate: data = ", data);
    const {toothnum, stateid, toothside, vid, pid} = data;
    (async () => {
  const rawResponse = await fetch(`/api/visits/${vid}/add`, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      VisitId: vid, 
      PatientId: pid,
      StateId: stateid,
      ToothNum: toothnum,
      Toothside: toothside})
  });
  const content = await rawResponse.json();
  console.log("Response from FETCH: ", content);
})();
  }
}
