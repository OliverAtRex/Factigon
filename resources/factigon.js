let form = document.getElementById("form");
let question = document.getElementById("question");
let answer = document.getElementById("answer");
form.onsubmit= function(e){
	fetch( '/ask?qu='+question.value)
		.then(function(response){
			response.text().then(function(ans){
				answer.textContent=ans;
			});
		});
	e.preventDefault();
};