"use strict";(()=>{window.addEventListener("load",()=>{let e=document.getElementById("country_selector"),t=document.getElementById("country_code");if(!e||!t)return;let c=e.style.color,o=t.innerText;function u(n){switch(n){case"AL":return"+355";case"LT":return"+370";case"LV":return"+371";case"EE":return"+372";case"MD":return"+373";case"RS":return"+381";case"ME":return"+382";case"XK":return"+383";case"BA":return"+387";case"MK":return"+389";case"LI":return"+423";default:return""}}function s(){let n=e.options[e.selectedIndex].disabled,r=u(e.options[e.selectedIndex].value);e.style.color=n?c:"black",t.style.color=r?"black":c,t.innerText=r||o}e.addEventListener("change",s),s()});})();
//# sourceMappingURL=bundle.js.map
