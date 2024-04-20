
<script> 
export default {
    props: ['src', 'uploader', 'date', 'likes','liked'], // some ID wont be visualized
    data: function(){
        return {
            liked: false,
            likes: 0,
           
        }
    },
    mounted(){
        alert(this.src)
    },
    
}
</script>


<!--- href for visualizing photo in another windows? what about -->
<template>
    <div>
        <div class="card">
            <img :src="src" class="card-img-top" />
            <div class="card-body">
                <p class="card-text">Uploader: {{ this.uploader }}</p> <!-- username href to profile -->
                <p class="card-text">Date: {{ this.date }}</p>
                <p class="card-text">Likes: {{ this.likes }}</p>
                <!-- this should become an icon probably -->
                <p class="card-text">Liked: {{ this.liked }}</p>
            </div>
        </div>
    </div>
</template>
