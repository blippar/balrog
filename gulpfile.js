var gulp = require('gulp');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var cleanCSS = require('gulp-clean-css');
var del = require('del');

var paths = {
  scripts: [
    'node_modules/jquery/dist/jquery.slim.js',
    'node_modules/popper.js/dist/umd/popper.js',
    'node_modules/bootstrap/dist/js/bootstrap.js',
    'node_modules/datatables/media/js/jquery.dataTables.js',
    'node_modules/datatables/media/js/dataTables.bootstrap4.js'
  ],
  styles: [
    'node_modules/bootstrap/dist/css/bootstrap.css',
    'node_modules/datatables/media/css/dataTables.bootstrap4.css',
    'node_modules/open-iconic/font/css/open-iconic-bootstrap.css'
  ],
  fonts: [
    'node_modules/open-iconic/font/fonts/*'
  ]
};

gulp.task('scripts', ['clean'], function() {
  return gulp.src(paths.scripts)
    .pipe(uglify())
    .pipe(concat('dist.js'))
    .pipe(gulp.dest('./dist/js/'));
});

gulp.task('styles', ['clean'], function() {
  return gulp.src(paths.styles)
    .pipe(concat('dist.css'))
    .pipe(cleanCSS())
    .pipe(gulp.dest('./dist/css/'));
});

gulp.task('fonts', ['clean'], function() {
  return gulp.src(paths.fonts)
    .pipe(gulp.dest('./dist/fonts/'));
});

gulp.task('clean', function() {
  return del('./dist/');
});

gulp.task('default', ['scripts','styles','fonts']);
