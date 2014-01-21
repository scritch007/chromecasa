function sendRequest(url, success, error){

	if (-1 != window.location.search.indexOf("debug=1")){
		url += "?debug=1"
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

function displayAlbums(albums){
	document.body.innerHTML = "";
	for(var i=0; i<albums.length; i++){
		var album = albums[i];
		var alb_div = document.createElement("div");
		var label = document.createElement("label");
		label.innerHTML = album.name;
		var img = document.createElement("img");
		img.src = album.icon;
		alb_div.appendChild(label);
		alb_div.appendChild(img);
		document.body.appendChild(alb_div);
		alb_div.onclick = function(){
			sendRequest("/album/" + album.id);
		}.bind(album);
	}
}

function main(){
	albums = sendRequest("album", displayAlbums);
}