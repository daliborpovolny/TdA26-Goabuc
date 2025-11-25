<script>
	let course_name = '';
	let course_description = '';

	function createCourse() {
		console.log(course_name, course_description);
	}
</script>

<div>
	<h1>Create new course</h1>

	<form on:submit={createCourse}>
		<label for="course_name">Name</label><br />
		<input type="text" bind:value={course_name} name="course_name" /><br />

		<label for="course_description">Description</label><br />
		<input type="text" bind:value={course_description} name="course_description" /><br />

		<button type="submit">Create</button>
	</form>
</div>
