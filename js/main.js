function sendRequest(url, success, error){

	//Check if page was loaded in debug mode
	url += window.location.search;

	var xhr = new XMLHttpRequest();
    xhr.open("GET", url);

    xhr.onload = function(e){
      var result = JSON.parse(xhr.responseText);
      success(result);
    }
    xhr.onerror = function(e){
      error();
    }
    //Not working because of cross site scripting
    try{
      xhr.send();
    }catch(err){

    }
}
function debug_mode(){
	return (-1 != window.location.search.indexOf("debug=1"));
}



function display_images(result){
	if (debug_mode() || (-1 != window.location.search.indexOf("noccast=1")))
	{
		//TODO DEBUG MODE...

		var slideshow_wrapper = document.getElementById("slideshow-wrapper");
		if (slideshow_wrapper === null){
			var display_overlay = document.createElement("div");
			display_overlay.className = "display_overlay";
			document.body.appendChild(display_overlay);
			var slideshow = document.createElement("div");
			slideshow.id = "slideshow";
			display_overlay.appendChild(slideshow);
			var slideshow_wrapper = document.createElement("div");
			slideshow_wrapper.id = "slideshow-wrapper";
			slideshow.appendChild(slideshow_wrapper)
		}

		SlideShow(result, slideshow_wrapper);
		return;
	}
	cast_api.sendMessage(cv_activity.activityId, 'BENJI', {images: result});
	cast_api.sendMessage(cv_activity.activityId, 'BENJI', {command: {name: 'loading', parameters:{load:false}}});
}
var g_albums = [];

function displayAlbums(albums){
	var main_display = document.getElementById("main_display");
	main_display.innerHTML = "";
	g_albums = albums;

	for(var i=0; i<albums.length; i++){
		var album = albums[i];
		var alb_div = document.createElement("div");
		alb_div.className = "album_div";
		var label = document.createElement("label");
		label.innerHTML = album.name;
		var img = document.createElement("img");
		img.src = album.icon;
		alb_div.appendChild(img);
		alb_div.appendChild(label);

		main_display.appendChild(alb_div);

		alb_div.onclick = function(){
			var album_id = album.id;
			return function(){
				if (! debug_mode() && !(-1 != window.location.search.indexOf("noccast=1"))) {
					cast_api.sendMessage(cv_activity.activityId, 'BENJI', {command: {name: 'loading', parameters:{load:true}}});
				}
				sendRequest("/album/" + album_id, display_images);
			}
		}();
	}
}

function main(){
	if(debug_mode())
	{
		//If we are in debug mode then we'll load the javascript file for the slideshow display
		var script = document.createElement("script");
		script.src = "/js/slideshow.js";
		document.body.appendChild(script);
	}
	albums = sendRequest("album", displayAlbums);
}