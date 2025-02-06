import{_ as w}from"./Card.vue_vue_type_script_setup_true_lang-PC315oxF.js";import{d,x as p,o as c,C as x,u as l,h as g,b as k,p as u,a as y,w as r,q as a,c as o,D as M,v as N,E as L,s as V}from"./index-CMT-UV0u.js";const A=d({__name:"CardContent",props:{class:{}},setup(s){const e=s;return(t,n)=>(c(),p("div",{class:x(l(g)("p-6 pt-0",e.class))},[k(t.$slots,"default")],2))}}),D=d({__name:"CardHeader",props:{class:{}},setup(s){const e=s;return(t,n)=>(c(),p("div",{class:x(l(g)("flex flex-col gap-y-1.5 p-6",e.class))},[k(t.$slots,"default")],2))}});/**
 * @license @tabler/icons-vue v3.29.0 - MIT
 *
 * This source code is licensed under the MIT license.
 * See the LICENSE file in the root directory of this source tree.
 */var v={outline:{xmlns:"http://www.w3.org/2000/svg",width:24,height:24,viewBox:"0 0 24 24",fill:"none",stroke:"currentColor","stroke-width":2,"stroke-linecap":"round","stroke-linejoin":"round"},filled:{xmlns:"http://www.w3.org/2000/svg",width:24,height:24,viewBox:"0 0 24 24",fill:"currentColor",stroke:"none"}};/**
 * @license @tabler/icons-vue v3.29.0 - MIT
 *
 * This source code is licensed under the MIT license.
 * See the LICENSE file in the root directory of this source tree.
 */const E=(s,e,t,n)=>({color:_="currentColor",size:h=24,stroke:C=2,title:m,class:P,...$},{attrs:B,slots:f})=>{let i=[...n.map(b=>u(...b)),...f.default?[f.default()]:[]];return m&&(i=[u("title",m),...i]),u("svg",{...v[s],width:h,height:h,...B,class:["tabler-icon",`tabler-icon-${e}`],...s==="filled"?{fill:_}:{"stroke-width":C??v[s]["stroke-width"],stroke:_},...$},i)};/**
 * @license @tabler/icons-vue v3.29.0 - MIT
 *
 * This source code is licensed under the MIT license.
 * See the LICENSE file in the root directory of this source tree.
 */var S=E("outline","git-branch","IconGitBranch",[["path",{d:"M7 18m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0",key:"svg-0"}],["path",{d:"M7 6m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0",key:"svg-1"}],["path",{d:"M17 6m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0",key:"svg-2"}],["path",{d:"M7 8l0 8",key:"svg-3"}],["path",{d:"M9 18h6a2 2 0 0 0 2 -2v-5",key:"svg-4"}],["path",{d:"M14 14l3 -3l3 3",key:"svg-5"}]]);const T={class:"flex items-center"},U={class:"p-2 rounded-full"},G={class:"ml-3"},I={class:"text-lg font-bold"},j=d({__name:"Leader",props:{leaderName:{}},setup(s){return(e,t)=>(c(),y(l(w),{class:"w-[250px] rounded-xl p-4 flex items-center"},{default:r(()=>[a("div",T,[a("div",U,[o(l(S),{class:"w-6 h-6"})]),a("div",G,[t[0]||(t[0]=a("p",{class:"text-gray-400 text-sm"},"Leader",-1)),a("p",I,M(e.leaderName),1)])])]),_:1}))}}),q={class:"flex"},H={class:"flex mr-[0.5em]"},K={class:"flex"},z=d({__name:"Dashboard",setup(s){const e=N("Unknown Leader");return L(()=>{e.value=window.TASKVAULT_LEADER||"Unknown Leader"}),(t,n)=>(c(),p("div",null,[o(l(w),null,{default:r(()=>[o(l(D),null,{default:r(()=>n[0]||(n[0]=[V("Welcome")])),_:1}),o(l(A),null,{default:r(()=>[a("div",q,[a("div",H,[a("div",K,[o(j,{"leader-name":e.value},null,8,["leader-name"])])])])]),_:1})]),_:1})]))}});export{z as default};
