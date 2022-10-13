const handleChange = (event)=> {
  const parent = event.target.closest('.card');
  const text = parent.querySelector('.js-text');
  console.log(text);
  const link = parent.querySelector('.js-title');
  link.classList.add('collapsed');
  text.classList.remove('collapsing');
  setTimeout(()=>{
      text.classList.remove('show');
  },300)
  
  console.log(link);
}