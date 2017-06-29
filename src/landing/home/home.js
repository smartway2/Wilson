$(document).ready(function(){
  console.log('heard')

  $('#prime').submit(function(e){
    e.preventDefault();
    console.log(this.elements[0].value);
    console.log(this.elements[1].value);
    const email = this.elements[0].value;
    const password = this.elements[1].value;
    $.post(`/login`, {email:email, password:password}, (data) => {
      console.log(data === 'continue');
      if(data === 'continue'){
        location.reload(true)
      } else {
        alert('Username or password invalid')
      }
    })
  })
})
