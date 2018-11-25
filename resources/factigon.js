let form = document.getElementById("form");
let question = document.getElementById("question");
let answer = document.getElementById("answer");

var thinking = false;

form.onsubmit= function(e){
	if (thinking)
		return;
	thinking = true;
	answer.textContent = "Thinking...";
	var qu = encodeURIComponent(question.value);
	fetch( '/ask?qu='+qu)
		.then(function(response){
			response.text().then(function(ans){
				answer.textContent=ans;
				question.select();
				thinking = false;
			});
		}).catch(function(err) {
			answer.textContent="ERROR: " + err;
			thinking = false;
		});
	e.preventDefault();
};