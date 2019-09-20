"use strict";

const tabletMin = 600;
const desktopMin = 900;

let focus = false;

function boxClick(event) {
	if (!focus) {
		event.target.closest(".box").classList.add("focus");	
		
	} 
	else {
		event.target.closest(".box").classList.remove("focus");
	}
}

document.getElementsByTagName("body")[0].addEventListener("transitionend", transitionEnd);
function transitionEnd(event) {
	if (event.propertyName == "height" && event.target.classList && 
			event.target.classList.contains("box")) {
		focus = !focus;
	}
}
