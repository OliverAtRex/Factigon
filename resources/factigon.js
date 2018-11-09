form.onsubmit= function(e){
  var ans=localStorage[question.value];
  answer.textContent=ans;
  new_answer.value=ans;
  e.preventDefault();
}
new_answer.onchange= function(e){
  localStorage[question.value]=new_answer.value;
}