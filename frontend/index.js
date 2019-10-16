//const URL = "http://localhost:8080/v1/convert"
const URL = "https://us-central1-yaml2go.cloudfunctions.net/yaml2go"

let go = document.getElementById("goGenerator")

let editor = ""

window.generatorCall=function (){
  let yamlData  = document.getElementById("yamlGenerator").value
  document.getElementById('yamlGenerator').style.border = "1px solid #ced4da"
  yamlData = editor.getValue()
  $.ajax({
    'url' : `${URL}`,
    'type' : 'POST',
    'data' : yamlData,
    'success' : function(data) { 
        document.getElementById("error").style.display="none" 
        document.getElementById("err-span").innerHTML="";     
        go.setValue(data)
    },
    'error' : function(jqXHR, request,error)
    {
      document.getElementById('yamlGenerator').style.border = "1px solid red"
      if (jqXHR.status == 400) {
        displayError('Invalid yaml format')
      } else {
        displayError('Something went wrong! Please report this to me@prasadg.dev')
      }
    }
  });

}

//Convert
document.getElementById("convert").addEventListener('click', ()=>{
   generatorCall()
})

//Clear YAML
document.getElementById('clearYaml').addEventListener('click',()=>{
  editor.setValue('')
})

//Clear Go
document.getElementById('clearGo').addEventListener('click',()=>{
  go.setValue('')
})


$(document).ready(function(){
    //code here...
    var input = $(".codemirror-textarea")[0];
    var output = $(".codemirror-textarea")[1];
    editor = CodeMirror.fromTextArea(input, {
        mode: "text/x-yaml",
    	lineNumbers : true
    });

    go = CodeMirror.fromTextArea(output, {
    	lineNumbers : true,
        mode: "text/x-go"
    });

    editor.setValue('# Paste your yaml here...\n'+
    'kind: test\n' +
    'metadata:\n' +
    '  name: cluster\n' +
    '  namespace: test-ns\n')

    go.setValue('// Yaml2Go\n'+
    'type Yaml2Go struct {\n' +
    '	Kind     string   `yaml:"kind"`\n' +
    '	Metadata Metadata `yaml:"metadata"`\n' +
    '}\n' +
    '\n' +
    '// Metadata\n' +
    'type Metadata struct {\n' +
    '	Name      string `yaml:"name"`\n' +
    '	Namespace string `yaml:"namespace"`\n' +
    '}\n')
});

function displayError(err){
  document.getElementById("err-span").innerHTML=err;
  document.getElementById("error").style.display="block"
}
