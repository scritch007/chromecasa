function start_warning(text){
	var elem = document.getElementById("warning-wrapper");

	elem.innerHTML = "";
	var label = document.createElement("label");
	label.id = "warning-text";
	label.innerHTML = text
	elem.appendChild(label);
	elem.style.opacity = 1;
	setTimeout(function(){elem.style.opacity=0}, 2000);
}

function Loader(){
	this.counter = 0;
}

Loader.prototype.start = function(){
	this.counter += 1;
	if (this.counter == 1){
		var elem =document.getElementById("loading");
		elem.style.display = "";
	}
}

Loader.prototype.stop = function(){
	this.counter -= 1;
	if (this.counter == 0){
		var elem =document.getElementById("loading");
		elem.style.display = "none";	
	}
}

loader = new Loader();