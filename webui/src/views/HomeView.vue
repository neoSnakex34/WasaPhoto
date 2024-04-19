<script>
export default {
	data: function() {
		return {
			username: localStorage.getItem('username'),
			userId: localStorage.getItem('userId'),
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {

			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();

			}
			this.loading = false;
		},
		
		doLogout() {
			localStorage.removeItem('userId');
			localStorage.removeItem('username');
			this.$router.push('/login');
		},
		

	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		
		<div class="d-flex r input-group pb-4 pt-4 border-bottom">
			<input type="text" id="findUser" v-model="findUser" class="form-control form-control-lg rounded" placeholder="search"/>
  			<div class="d-flex justify-content-between">
				<input class="form-control form-control-lg ms-3 me-2" type="file" id="formFile">
				<button class="btn btn-primary btn-lg rounded-pill fw-bold">Upload</button>
			</div>
		</div>

		<div class="d-flex justify-content-between align-items-center pt-2 pb-2  border-bottom">
			<h1 class="h2"><strong>{{ this.username }}</strong>'s feed</h1>
				<div>
					<button type="button" class="btn btn-success btn-lg fw-bold rounded-pill" @click="refresh">
						Refresh
					</button>
					
				</div>
			
			
		</div>


		<!-- <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg> -->
	</div>
</template>

<style>
</style>
