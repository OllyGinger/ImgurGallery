document.getElementById("rte").focus();
document.body.addEventListener("paste", function(e) {
	for (var i = 0; i < e.clipboardData.items.length; i++) {
		
		if (e.clipboardData.items[i].kind == "file" && e.clipboardData.items[i].type.match(/image.*/) ) {
			var imageFile = e.clipboardData.items[i].getAsFile();
			upload( imageFile );
		}
	}
});

document.body.addEventListener("dragover", function(e) {
    e.stopPropagation();
    e.preventDefault();
}, false);

document.body.addEventListener("drop", function(e) {
    e.stopPropagation();
    e.preventDefault();

	for( var i = 0; i < e.dataTransfer.files.length; i++ ) {
		upload( e.dataTransfer.files[i] );
	}
	
}, false );

function upload(file) {
	if (!file || !file.type.match(/image.*/)) return;

	var fd = new FormData();
	fd.append("image", file); 
	fd.append("key", imgurAPIKey);
	
	document.getElementById("status").innerHTML = "UPLOADING...";
	
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "http://api.imgur.com/2/upload.json");
	xhr.onload = function() {
		var data = JSON.parse(xhr.responseText);

		// Upload the values to our own DB
		var store = new FormData();
		store.append( "hash", data.upload.image.hash );
		store.append( "deletehash", data.upload.image.deletehash );
		store.append( "orig", data.upload.links.original.split('/').pop() );
		store.append( "thumb", data.upload.links.small_square.split('/').pop() );

		var xhr2 = new XMLHttpRequest();
		xhr2.open( "POST", document.location + "uploaded" );
		xhr2.onload = function() {
			console.log( xhr2.responseText );
			window.location.reload()
		}
		xhr2.send( store );
	}
   
   xhr.send(fd);
}
