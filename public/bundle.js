"use strict";(()=>{document.body.addEventListener("htmx:beforeSwap",e=>{if(e.detail.pathInfo.requestPath==="/htmx/search"){switch(e.detail.xhr.status){case 404:case 500:break;default:return}e.detail.isError=!1,e.detail.shouldSwap=!0}});document.body.addEventListener("htmx:afterSwap",e=>{if(e.detail.pathInfo.requestPath!=="/htmx/search")return;let t=document.getElementById("search"),r=t==null?void 0:t.nextElementSibling;if(!t||!r)return;function n(){r.remove(),t.removeEventListener("input",n)}t.addEventListener("input",n)});document.body.addEventListener("htmx:afterSwap",e=>{if(!e.detail.failed){switch(e.detail.pathInfo.requestPath){case"/htmx/checkout":case"/htmx/address/validate":break;default:return}d()}});function d(){let e=document.getElementById("country_selector"),t=document.getElementById("country_code");if(!e||!t)return;let r=e.style.color,n=t.innerText;function a(){let s=e.options[e.selectedIndex].disabled,o=i(e.options[e.selectedIndex].value);e.style.color=s?r:"black",t.style.color=o?"black":r,t.innerText=o||n}e.addEventListener("change",a),a()}function i(e){switch(e){case"AL":return"+355";case"LT":return"+370";case"LV":return"+371";case"EE":return"+372";case"MD":return"+373";case"RS":return"+381";case"ME":return"+382";case"XK":return"+383";case"BA":return"+387";case"MK":return"+389";case"LI":return"+423";default:return""}}d();document.body.addEventListener("htmx:beforeSwap",e=>{if(e.detail.pathInfo.requestPath.startsWith("/htmx/product")){switch(e.detail.xhr.status){case 404:break;default:return}e.detail.isError=!1,e.detail.shouldSwap=!0}});document.body.addEventListener("htmx:beforeSwap",e=>{if(e.detail.isError||!e.detail.pathInfo.requestPath.startsWith("/htmx/product")||e.detail.serverResponse!="")return;let t=document.getElementById("placeholder");t&&(l(t),e.detail.shouldSwap=!1)});document.body.addEventListener("htmx:afterSwap",e=>{if(e.detail.failed||!e.detail.pathInfo.requestPath.startsWith("/htmx/add_to_bag"))return;let t=document.getElementById("placeholder");if(!t)return;let r=t.getElementsByTagName("img");if(r.length==0)return;let n=0;for(let a of r)a.onload=()=>{if(r.length!=++n)return;let s=document.getElementById("placeholder-close");t.style.transform="translateY(0%)",t.onclick=u(setTimeout(()=>l(t),2e3)),s&&(s.onclick=()=>l(t))}});function u(e){return t=>{t.preventDefault(),clearTimeout(e)}}function l(e){e.style.transform="",setTimeout(()=>{let t=e.querySelector("ul");t&&(t.innerHTML="")},200)}})();
//# sourceMappingURL=bundle.js.map
