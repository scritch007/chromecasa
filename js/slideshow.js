var display_timeout = 3000;

//TODO check with the actual display screen information
var display_height = window.screen.availHeight;
var display_width = window.screen.availWidth;

function SlideShow(images, in_display){
	this.current_image = 0;
	this.images = images;
	this.display = in_display;
	this.timeout = null;
	this.image_list = [];
	this.offset = Math.min(3, images.length);
}

var current_slideshow = null;

SlideShow.prototype.pause = function(){
	if (this.timeout)
	{
		window.clearTimeout(this.timeout);
		this.timeout = null;
	}
}
SlideShow.prototype.add = function(src){
	var img = document.createElement("img");
	var img_height = Number(src.height);
	var img_width = Number(src.width);
	var width_ratio = 1;
	var heigh_ratio = 1;
	if (img_width > display_width){
		width_ratio = display_width/img_width;
	}
	if (img_height > display_height){
		heigh_ratio = display_height/img_height;
	}
	var ratio = Math.min(width_ratio, heigh_ratio);
	img.src = src.content;
	img.width = ratio * img_width;
	img.height = ratio * img_height;
	img.className = "hide";
	this.display.appendChild(img);
	this.image_list.push(img);
	return img;
}
SlideShow.prototype.start = function(){
	this.display.innerHTML = "";
	if (current_slideshow){
		current_slideshow.pause();
	}
	current_slideshow = this;
	//PRELOAD some images

	
	for(var j = 0; j < this.offset; j++){
		img = this.add(this.images[j])
	}
	var slideshow = this;

	this.display_me();
}
SlideShow.prototype.display_me = function(){
		
	var i = this.current_image;
	if (i < this.images.length){

		this.next();
		var slideshow = this;
		this.timeout = setTimeout(function(){ 
			 
			return function(){
				slideshow.display_me()
			}
		}(), 3000);
	}else{
		//display.parentNode.removeChild(display);
		this.image_list[i-1].className = "hide";
	}
}

SlideShow.prototype.previous = function(){
	//Stop the current timer
	var slideshow = this;
	this.pause();
	//Hide current image
	if (this.current_image != 0)
	{
		var current_image = this.image_list[this.current_image - 1];

		current_image.className = "hide";
	}
	this.current_image -= 1;
	this.image_list[this.current_image - 1].className = "show";
}

SlideShow.prototype.next = function(){
	//clear running timer.
	this.pause();
	var i = this.current_image;
	if (this.images.length <= i){
		return
	}
	if (i < (this.images.length - this.offset))
	{
		//We want to preload some other images so start loading the images
		img = this.add(this.images[i])
	}
	if (0 != i)
	{
		this.image_list[i-1].className ="hide";
	}
	this.image_list[i].className = "show";
	i += 1;
	this.current_image = i;
}