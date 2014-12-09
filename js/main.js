function sendRequest(url, success, error){

	//Check if page was loaded in debug mode
	if (-1 != url.indexOf("?"))
	{
		url += "&"+window.location.search.substr(1);
	}else{
		url += window.location.search;
	}

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
	return (-1 != window.location.search.indexOf("provider=debug"));
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

function addAlbums(albums){
	g_albums.push(albums);
	displayAlbums(albums);
	document.getElementById("up_button").style.display = "";
}



function displayAlbums(albums){
	var main_display = document.getElementById("main_display");
	main_display.innerHTML = "";

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

		var display_label = document.createElement("label");
		display_label.innerHTML = "display "+album.display;
		alb_div.appendChild(display_label);
		var browse_label = document.createElement("label");
		browse_label.innerHTML = "browse "+album.browse;
		alb_div.appendChild(browse_label);

		main_display.appendChild(alb_div);

		if (album.display || album.browse)
		{
			var action_div = document.createElement("div");
			alb_div.appendChild(action_div);

			action_div.className = "action_div hidden";
			alb_div.addEventListener("mouseover", function(){
				this.className = "action_div";
			}.bind(action_div)
			);
			alb_div.addEventListener("mouseout", function(){
				this.className = "action_div hidden";
			}.bind(action_div));

			if (album.display)
			{
				var play = document.createElement("label");
				play.className = "action_element";
				play.innerHTML = "play";
				play.addEventListener("click",
					function(){
						if (! debug_mode() && !(-1 != window.location.search.indexOf("noccast=1"))) {
							cast_api.sendMessage(cv_activity.activityId, 'BENJI', {command: {name: 'loading', parameters:{load:true}}});
						}
						sendRequest("/display/" + this.id, display_images);
					}.bind(album)
				);
				action_div.appendChild(play);
			}
			if(album.browse){
				var browse = document.createElement("label");
				browse.className = "action_element";
				browse.innerHTML = "browse";
				browse.addEventListener("click",
					function(){
						sendRequest("/browse?path=" + this.id, addAlbums);
					}.bind(album)
				);
				action_div.appendChild(browse);
			}
		}
	}
}

function main(){
	document.getElementById("up_button").addEventListener("click",
		function(){
			if (g_albums.length > 0){
				g_albums = g_albums.slice(0, g_albums.length-1);
			}
			displayAlbums(g_albums[g_albums.length-1]);
			if (1 == g_albums.length){
				document.getElementById("up_button").style.display = "none";
			}
		}
	)
	if(debug_mode())
	{
		//If we are in debug mode then we'll load the javascript file for the slideshow display
		var script = document.createElement("script");
		script.src = "/js/slideshow.js";
		document.body.appendChild(script);
	}
	albums = sendRequest("browse?path=/", addAlbums);
}