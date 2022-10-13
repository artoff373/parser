// const handleChange = (event)=> {
//   const parent = event.target.closest('.card');
//   const text = parent.querySelector('.js-text');
//   const link = parent.querySelector('.js-title');
//   link.classList.add('collapsed');
//   text.classList.remove('collapsing');
//   text.classList.remove('show');
// }

const handleChange = (event) => {
  const parent = event.target.closest('.card');
  const text = $('.js-text', parent);
  text.collapse('hide');
}