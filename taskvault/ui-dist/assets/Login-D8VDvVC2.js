import{d as m,e as V,a as U,o as p,w as g,b as x,m as C,u as a,j as h,R as I,C as S,D as $,x as _,E as w,F as N,p as u,t as T,G as j,q as l,c as d,H as B,s as M,B as R,U as E,_ as L}from"./index-cQHmgfoe.js";import{b as P,a as A,_ as q}from"./index-Dm4-eGFp.js";const z=m({__name:"ToastAction",props:{altText:{},asChild:{type:Boolean},as:{},class:{}},setup(s){const o=s,t=V(()=>{const{class:e,...n}=o;return n});return(e,n)=>(p(),U(a(I),C(t.value,{class:a(h)("inline-flex h-8 shrink-0 items-center justify-center rounded-md border bg-transparent px-3 text-sm font-medium transition-colors hover:bg-secondary focus:outline-none focus:ring-1 focus:ring-ring disabled:pointer-events-none disabled:opacity-50 group-[.destructive]:border-muted/40 group-[.destructive]:hover:border-destructive/30 group-[.destructive]:hover:bg-destructive group-[.destructive]:hover:text-destructive-foreground group-[.destructive]:focus:ring-destructive",o.class)}),{default:g(()=>[x(e.$slots,"default")]),_:3},16,["class"]))}}),b=m({__name:"Input",props:{defaultValue:{},modelValue:{},class:{}},emits:["update:modelValue"],setup(s,{emit:o}){const t=s,n=P(t,"modelValue",o,{passive:!0,defaultValue:t.defaultValue});return(f,r)=>S((p(),_("input",{"onUpdate:modelValue":r[0]||(r[0]=i=>N(n)?n.value=i:null),class:w(a(h)("flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",t.class))},null,2)),[[$,a(n)]])}}),D=m({__name:"Card",props:{class:{}},setup(s){const o=s;return(t,e)=>(p(),_("div",{class:w(a(h)("rounded-xl border bg-card text-card-foreground shadow",o.class))},[x(t.$slots,"default")],2))}});/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const G=s=>s.replace(/([a-z0-9])([A-Z])/g,"$1-$2").toLowerCase();/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */var c={xmlns:"http://www.w3.org/2000/svg",width:24,height:24,viewBox:"0 0 24 24",fill:"none",stroke:"currentColor","stroke-width":2,"stroke-linecap":"round","stroke-linejoin":"round"};/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const H=({size:s,strokeWidth:o=2,absoluteStrokeWidth:t,color:e,iconNode:n,name:f,class:r,...i},{slots:v})=>u("svg",{...c,width:s||c.width,height:s||c.height,stroke:e||c.stroke,"stroke-width":t?Number(o)*24/Number(s):o,class:["lucide",`lucide-${G(f??"icon")}`],...i},[...n.map(k=>u(...k)),...v.default?[v.default()]:[]]);/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const y=(s,o)=>(t,{slots:e})=>u(H,{...t,iconNode:o,name:s},e);/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const K=y("RectangleEllipsisIcon",[["rect",{width:"20",height:"12",x:"2",y:"6",rx:"2",key:"9lu3g6"}],["path",{d:"M12 12h.01",key:"1mp3jc"}],["path",{d:"M17 12h.01",key:"1m0b6t"}],["path",{d:"M7 12h.01",key:"eqddd0"}]]);/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const F=y("UserIcon",[["path",{d:"M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2",key:"975kel"}],["circle",{cx:"12",cy:"7",r:"4",key:"17ys0d"}]]),J={class:"h-screen login-bg bg-slate-50"},O={class:"flex h-full justify-center items-center"},Y={class:"h-max min-w-[16rem] w-1/4 max-w-[24rem] text-center"},Z={class:"mb-3 relative w-full max-w-sm items-center"},Q={class:"absolute start-0 inset-y-0 flex items-center justify-center px-2"},W={class:"mb-3 relative w-full max-w-sm items-center"},X={class:"absolute start-0 inset-y-0 flex items-center justify-center px-2"},ee=m({__name:"Login",setup(s){const o=R(),{toast:t}=T(),e=j({Username:"",Password:""}),n=()=>{e.Username==="admin"&&e.Password==="admin"?(localStorage.setItem(E,JSON.stringify(e)),t({description:"Login successfully."}),o.push("/")):t({title:"Uh oh! Something went wrong.",description:"Wrong username or password.",variant:"destructive",action:u(z,{altText:"Try again"},{default:()=>"Try again"})})};return(f,r)=>(p(),_("div",J,[l("div",O,[l("div",Y,[r[3]||(r[3]=l("div",{class:"inline-flex mt-4 mb-8 items-center"},[l("img",{src:A,class:"h-12 mr-2"}),l("h1",{class:"font-bold text-4xl font-mono"},"Syncra")],-1)),d(a(D),{class:"p-6 shadow-lg"},{default:g(()=>[l("form",{onSubmit:B(n,["prevent"])},[l("div",Z,[d(a(b),{id:"user",modelValue:e.Username,"onUpdate:modelValue":r[0]||(r[0]=i=>e.Username=i),class:"pl-10 w-full mt-1",placeholder:"admin"},null,8,["modelValue"]),l("span",Q,[d(a(F),{class:"size-6 text-muted-foreground"})])]),l("div",W,[d(a(b),{id:"password",modelValue:e.Password,"onUpdate:modelValue":r[1]||(r[1]=i=>e.Password=i),type:"password",class:"pl-10 w-full mt-1",placeholder:""},null,8,["modelValue"]),l("span",X,[d(a(K),{class:"size-6 text-muted-foreground"})])]),d(a(q),{type:"submit",class:"w-full mt-3"},{default:g(()=>r[2]||(r[2]=[M("SIGN IN")])),_:1})],32)]),_:1})])])]))}}),oe=L(ee,[["__scopeId","data-v-403b28c6"]]);export{oe as default};
