function downURL(output,url){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        // data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = data;
        // document.getElementById(output).innerHTML = data;
      }
    };
    req.open("POST","/download",true);
    req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    var str = "url="+url
    req.send(str);
}

function pushMp3Tag(jsondata,output,id){
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      data = JSON.parse(data);
      console.log(data);
      getMp3Tag(output,id)
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("POST","/edit/"+id,true);
  req.send(jsondata);

}

function editMp3Tag(output,id) {
  document.getElementById("lyrics").innerHTML = document.getElementById("lyrics").innerHTML.replace(/<div><br><\/div>/g, '<div></div>')
  var lyrics = document.getElementById("lyrics").innerHTML.replace(/<div>/g, '').replace(/<\/div>/g, '\r\n').replace(/<br>/g,'\r\n')
  // var lyrics = document.getElementById("lyrics").innerText.replace(/\n/g, '\r\n')
  var json = {
    "Title":document.getElementById("title").innerHTML,
    "Artist":document.getElementById("artist").innerHTML,
    "Album":document.getElementById("album").innerHTML,
    "Year":document.getElementById("year").innerHTML,
    "Lyrics":lyrics,
  }
  // console.log(json)
  var jsonstr =JSON.stringify(json)
  console.log(jsonstr)
  pushMp3Tag(jsonstr,output,id)
}

function getMp3List(output) {

  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      data = JSON.parse(data);
      console.log(data)
      document.getElementById(output).innerHTML = createDownloadLink(data);
    }
  };
  req.open("GET","/list",true);
  req.send( null );
}

function getMp3Tag(output,id) {
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      data = JSON.parse(data);
      console.log(data)
      document.getElementById(output).innerHTML = createMp3Tag(output,data,id);
    }
  };
  req.open("GET","/view/"+id,true);
  req.send( null );
}

function createMp3Tag(outputid,data,id){
  var output = "";
  output += "<div><button onclick=\""+"editMp3Tag('"+outputid+"'"+",'"+id+"')"+"\">write</button></div>"
  output += createTextDiv(data.Title,"title")
  output += createTextDiv(data.Artist,"artist")
  output += createTextDiv(data.Album,"album")
  output += createTextDiv(data.Year,"year")
  output += createTextDiv(data.Lyrics.replace(/\r\n/g, '<br>'),"lyrics")
  return output
}

function createTextDiv(str,id) {
  var output =""
  if (id=="lyrics"){
    output += "<div class=\"lyricsrow\">"+id+"</div>"
  }else{
    output += "<div class=\"row\">"+id+"</div>"
  }
  output += "<div class id=\""+id+"\" contenteditable>"+str+"</div>"
  return output
}

function createDownloadLink(data) {
  var output = ""
  for(var i=0;i<data.length;i++) {
    var tmp = data[i]
    if (i>0) {
      output += "<br>"
    }
    output += "<a href='/mp3/"+tmp.No+"'>"+tmp.Name+"<a>"+" "+"<audio controls src=\"/mp3/"+tmp.No+"\" type=\"audio/mp3\">"+"</audio>"
    output += "<button onclick=\""+"getMp3Tag('edit','"+tmp.No+"')"+"\">view</button>"
  }
  return output
}

var health_down_flag = false;
var health_down_url =[];
var health_down_cmd ="";

function healthck(output) {
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      data = JSON.parse(data);
      health_down_flag = data.Downflag
      health_down_url = data.DownUrl
      health_down_cmd = data.DownCmd
      if (health_down_flag) {
        document.getElementById(output).innerHTML = health_down_cmd;
      }else{
        document.getElementById(output).innerHTML = "not down load";
        
      }
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("GET","/health",true);
  req.send( null );
}

var timer = null;

function watch_cmd(output) {
  if (timer != null)  {
    clearInterval(timer);
  }else{
    console.log( "timerstart");
  }
  healthck(output)
  if (health_down_flag) {
    if ((health_down_cmd != "")&&(health_down_cmd != null)){
      document.getElementById(output).innerHTML = health_down_cmd;
    }
    timer = setTimeout(watch_cmd,2000,output)
  }else {
    document.getElementById(output).innerHTML = "Check End";
    console.log("TimeEnd");
    timer = null;
  }
}