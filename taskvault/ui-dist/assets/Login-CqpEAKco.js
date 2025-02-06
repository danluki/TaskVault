import{d as f,e as y,a as k,o as g,w as p,b as V,m as U,u as l,h as _,R as I,F as C,G as T,x as b,C as $,H as B,p as u,t as M,I as N,q as n,c as d,J as P,s as j,B as R,K as L,_ as S}from"./index-CMT-UV0u.js";import{b as A,a as E,_ as q}from"./index-ChpK9YOp.js";import{_ as z}from"./Card.vue_vue_type_script_setup_true_lang-PC315oxF.js";const G=f({__name:"ToastAction",props:{altText:{},asChild:{type:Boolean},as:{},class:{}},setup(s){const a=s,t=y(()=>{const{class:e,...r}=a;return r});return(e,r)=>(g(),k(l(I),U(t.value,{class:l(_)("inline-flex h-8 shrink-0 items-center justify-center rounded-md border bg-transparent px-3 text-sm font-medium transition-colors hover:bg-secondary focus:outline-none focus:ring-1 focus:ring-ring disabled:pointer-events-none disabled:opacity-50 group-[.destructive]:border-muted/40 group-[.destructive]:hover:border-destructive/30 group-[.destructive]:hover:bg-destructive group-[.destructive]:hover:text-destructive-foreground group-[.destructive]:focus:ring-destructive",a.class)}),{default:p(()=>[V(e.$slots,"default")]),_:3},16,["class"]))}}),v=f({__name:"Input",props:{defaultValue:{},modelValue:{},class:{}},emits:["update:modelValue"],setup(s,{emit:a}){const t=s,r=A(t,"modelValue",a,{passive:!0,defaultValue:t.defaultValue});return(m,o)=>C((g(),b("input",{"onUpdate:modelValue":o[0]||(o[0]=i=>B(r)?r.value=i:null),class:$(l(_)("flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",t.class))},null,2)),[[T,l(r)]])}});/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const H=s=>s.replace(/([a-z0-9])([A-Z])/g,"$1-$2").toLowerCase();/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */var c={xmlns:"http://www.w3.org/2000/svg",width:24,height:24,viewBox:"0 0 24 24",fill:"none",stroke:"currentColor","stroke-width":2,"stroke-linecap":"round","stroke-linejoin":"round"};/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const K=({size:s,strokeWidth:a=2,absoluteStrokeWidth:t,color:e,iconNode:r,name:m,class:o,...i},{slots:h})=>u("svg",{...c,width:s||c.width,height:s||c.height,stroke:e||c.stroke,"stroke-width":t?Number(a)*24/Number(s):a,class:["lucide",`lucide-${H(m??"icon")}`],...i},[...r.map(x=>u(...x)),...h.default?[h.default()]:[]]);/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const w=(s,a)=>(t,{slots:e})=>u(K,{...t,iconNode:a,name:s},e);/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const D=w("RectangleEllipsisIcon",[["rect",{width:"20",height:"12",x:"2",y:"6",rx:"2",key:"9lu3g6"}],["path",{d:"M12 12h.01",key:"1mp3jc"}],["path",{d:"M17 12h.01",key:"1m0b6t"}],["path",{d:"M7 12h.01",key:"eqddd0"}]]);/**
 * @license lucide-vue-next v0.474.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const F=w("UserIcon",[["path",{d:"M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2",key:"975kel"}],["circle",{cx:"12",cy:"7",r:"4",key:"17ys0d"}]]),J={class:"h-screen login-bg bg-slate-50"},Z={class:"flex h-full justify-center items-center"},O={class:"h-max min-w-[16rem] w-1/4 max-w-[24rem] text-center"},Q={class:"mb-3 relative w-full max-w-sm items-center"},W={class:"absolute start-0 inset-y-0 flex items-center justify-center px-2"},X={class:"mb-3 relative w-full max-w-sm items-center"},Y={class:"absolute start-0 inset-y-0 flex items-center justify-center px-2"},ee=f({__name:"Login",setup(s){const a=R(),{toast:t}=M(),e=N({Username:"",Password:""}),r=()=>{e.Username==="admin"&&e.Password==="admin"?(L({Username:"admin",Password:"admin"}),t({description:"Login successfully."}),a.push("/")):t({title:"Uh oh! Something went wrong.",description:"Wrong username or password.",variant:"destructive",action:u(G,{altText:"Try again"},{default:()=>"Try again"})})};return(m,o)=>(g(),b("div",J,[n("div",Z,[n("div",O,[o[3]||(o[3]=n("div",{class:"inline-flex mt-4 mb-8 items-center"},[n("img",{src:E,class:"h-12 mr-2"}),n("h1",{class:"font-bold text-4xl font-mono"},"Syncra")],-1)),d(l(z),{class:"p-6 shadow-lg"},{default:p(()=>[n("form",{onSubmit:P(r,["prevent"])},[n("div",Q,[d(l(v),{id:"user",modelValue:e.Username,"onUpdate:modelValue":o[0]||(o[0]=i=>e.Username=i),class:"pl-10 w-full mt-1",placeholder:"admin"},null,8,["modelValue"]),n("span",W,[d(l(F),{class:"size-6 text-muted-foreground"})])]),n("div",X,[d(l(v),{id:"password",modelValue:e.Password,"onUpdate:modelValue":o[1]||(o[1]=i=>e.Password=i),type:"password",class:"pl-10 w-full mt-1",placeholder:""},null,8,["modelValue"]),n("span",Y,[d(l(D),{class:"size-6 text-muted-foreground"})])]),d(l(q),{type:"submit",class:"w-full mt-3"},{default:p(()=>o[2]||(o[2]=[j("SIGN IN")])),_:1})],32)]),_:1})])])]))}}),ae=S(ee,[["__scopeId","data-v-f8ac701d"]]);export{ae as default};
