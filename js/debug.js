function start(){
	var result = [
		{
			name: "Photo 1",
			content: "/img/photo1.jpg",
			height:"384",
			width: "512"
		},
		{
			name: "Photo 2",
			content: "/img/default.png",
			height: "59",
			width: "120"
		},
		{
			name: "Photo 3",
			content: "/img/photo3.jpeg",
			height: "3264",
			width: "2448"
		}
	];
	for (var i=1; i < 22; i++){
		result.push({
			name: "Photo " + (i+3),
			content: "/img/200." + i + ".png",
			height: "200",
			width:"400"
		});
	}

	var slideshow_wrapper = document.getElementById("slideshow-wrapper");
	var slideshow = new SlideShow(result, slideshow_wrapper)
	slideshow.start();
}
function previous(){
	current_slideshow.previous();
}

function next(){
	current_slideshow.next();
}

function pause(){
	current_slideshow.pause();
}
function continue_(){
	current_slideshow.display_me();
}

function warning(){
	start_warning("Warning message exemple");
}

function loading(){
	loader.start();
	setTimeout(function(){loader.stop()}, 4000);
}