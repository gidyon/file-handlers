(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["pathogens"],{"2fdb":function(t,e,a){"use strict";var n=a("5ca1"),r=a("d2c8"),i="includes";n(n.P+n.F*a("5147")(i),"String",{includes:function(t){return!!~r(this,t,i).indexOf(t,arguments.length>1?arguments[1]:void 0)}})},5147:function(t,e,a){var n=a("2b4c")("match");t.exports=function(t){var e=/./;try{"/./"[t](e)}catch(a){try{return e[n]=!1,!"/./"[t](e)}catch(r){}}return!0}},6762:function(t,e,a){"use strict";var n=a("5ca1"),r=a("c366")(!0);n(n.P,"Array",{includes:function(t){return r(this,t,arguments.length>1?arguments[1]:void 0)}}),a("9c6c")("includes")},aae3:function(t,e,a){var n=a("d3f4"),r=a("2d95"),i=a("2b4c")("match");t.exports=function(t){var e;return n(t)&&(void 0!==(e=t[i])?!!e:"RegExp"==r(t))}},b7e3:function(t,e,a){"use strict";a.r(e);var n=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("v-container",{attrs:{fluid:"","pt-0":""}},[a("v-layout",{attrs:{"justify-center":""}},[a("v-flex",{attrs:{xs12:"",sm12:"",md9:"",lg6:""}},[a("v-layout",{attrs:{column:""}},[a("v-flex",{attrs:{"mb-3":"",bw:""}},[a("v-text-field",{attrs:{label:"Search viruses, bacteria, protozoa, helminths",placeholder:"Start typing to search","append-icon":t.inputIcon,outline:""},on:{"click:append":function(e){t.searchPathogen=""}},model:{value:t.searchPathogen,callback:function(e){t.searchPathogen=e},expression:"searchPathogen"}})],1),t._l(t.filteredPathogens,function(e){return a("v-flex",{key:e.pathogenId,attrs:{"mb-3":"","pa-2":"",card:"",pcoh:""},on:{click:function(a){return t.goto(e.pathogenId)}}},[a("v-layout",{attrs:{"justify-space-between":"","align-center":""}},[a("v-avatar",{attrs:{size:"24",color:"grey lighten-4"}},[a("img",{attrs:{src:t.pathogenUrl(e.type),alt:"avatar"}})]),a("v-flex",{attrs:{xs12:"","pl-3":""}},[a("span",{staticClass:"subheading font-weight-light"},[t._v(t._s(e.pathogenName))])]),a("v-flex",{attrs:{"pl-4":""}},[a("v-icon",[t._v("keyboard_arrow_right")])],1)],1)],1)})],2)],1)],1)],1)},r=[],i=(a("6762"),a("2fdb"),a("cebc")),c=a("2f62"),o={created:function(){this.getPathogens(1)},computed:Object(i["a"])({},Object(c["c"])(["pathogens"]),{filteredPathogens:function(){var t=this;return this.pathogens.filter(function(e){return e.pathogenName.toLowerCase().includes(t.searchPathogen.toLowerCase())})},inputIcon:function(){return""!==this.searchPathogen?"clear":"search"}}),data:function(){return{searchPathogen:""}},methods:Object(i["a"])({},Object(c["b"])(["getPathogens"]),{goto:function(t){this.$router.push("/pathogens/".concat(t))},pathogenUrl:function(t){switch(t){case"virus":return"https://img.scoop.it/6Tvnpt9ZkSop-G9oSN5Objl72eJkfbmt4t8yenImKBVvK0kTmF0xjctABnaLJIm9";case"fungi":return"https://upload.wikimedia.org/wikipedia/commons/thumb/3/34/Candida_pap_1.jpg/220px-Candida_pap_1.jpg";case"protozoa":return"https://katekearapelen.files.wordpress.com/2010/03/1760532740_20189a8cdc.jpg";case"rickettsia":return"https://upload.wikimedia.org/wikipedia/commons/thumb/8/86/Rickettsia_rickettsii.jpg/220px-Rickettsia_rickettsii.jpg";case"helminth":return"https://www.infectiousdisease.cam.ac.uk/images/Schistomsoma%20mansoni.jpg/image_mini";default:return"https://previews.123rf.com/images/frenta/frenta1611/frenta161100123/66972770-pathogen-bacteria-on-the-surface-3d-render.jpg"}}})},s=o,u=a("2877"),p=Object(u["a"])(s,n,r,!1,null,null,null);e["default"]=p.exports},d2c8:function(t,e,a){var n=a("aae3"),r=a("be13");t.exports=function(t,e,a){if(n(e))throw TypeError("String#"+a+" doesn't accept regex!");return String(r(t))}}}]);
//# sourceMappingURL=pathogens.3cd7ae23.js.map