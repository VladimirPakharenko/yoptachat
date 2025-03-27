var btnup = document.querySelector('.btnup')
var btnin = document.querySelector('.btnin')
var formreg = document.querySelector('.form_registration')
var formauth = document.querySelector('.form_authorization')
var a1 = document.querySelector('.sign_up')
var a2 = document.querySelector('.sign_in')

a2.style.borderBottom = '1px solid var(--fon)'
btnup.addEventListener('click', (e)=>{
    formreg.classList.remove('hidden')  
    formauth.classList.add('hidden')  
    a1.style.borderBottom = ''
    a2.style.borderBottom = '1px solid var(--fon)'
})
btnin.addEventListener('click', (e)=>{
    formauth.classList.remove('hidden')  
    formreg.classList.add('hidden')  
    a2.style.borderBottom = ''
    a1.style.borderBottom = '1px solid var(--fon)'
})