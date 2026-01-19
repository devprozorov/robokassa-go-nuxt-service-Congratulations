<template>
  <div class="absolute inset-0 pointer-events-none overflow-hidden">
    <!-- vignette + glow dots -->
    <div class="absolute inset-0 vignette"></div>
    <div class="absolute -bottom-24 left-0 right-0 h-40 glow"></div>

    <!-- Garland (newyear) -->
    <div v-if="kind === 'newyear'" class="absolute top-0 left-0 right-0 h-10 garland"></div>

    <!-- Snowflakes -->
    <template v-if="mounted && kind === 'newyear'">
      <span
        v-for="(f, i) in flakes"
        :key="i"
        class="snowflake"
        :style="{
          left: f.left + '%',
          width: f.size + 'px',
          height: f.size + 'px',
          opacity: f.opacity,
          animationDuration: f.duration + 's',
          animationDelay: f.delay + 's'
        }"
      />
    </template>

    <!-- Balloons + sprinkles -->
    <template v-if="mounted && kind === 'birthday'">
      <span v-for="(b, i) in balloons" :key="i" class="balloon"
        :style="{ left: b.left + '%', animationDuration: b.duration + 's', animationDelay: b.delay + 's' }"
      />
      <div class="absolute inset-0 sprinkle"></div>
    </template>

    <!-- Hearts -->
    <template v-if="mounted && (kind === 'proposal' || kind === 'anniversary')">
      <span v-for="(h, i) in hearts" :key="i" class="heart"
        :style="{ left: h.left + '%', animationDuration: h.duration + 's', animationDelay: h.delay + 's', opacity: h.opacity }"
      >❤</span>
    </template>

    <!-- Sparks -->
    <template v-if="mounted && kind === 'promotion'">
      <span v-for="(s, i) in sparks" :key="i" class="spark"
        :style="{ left: s.left + '%', top: s.top + '%', animationDelay: s.delay + 's' }"
      />
    </template>

    <!-- Stars -->
    <template v-if="mounted && (kind === 'baby' || kind === 'religious')">
      <span v-for="(s, i) in stars" :key="i" class="star"
        :style="{ left: s.left + '%', top: s.top + '%', animationDelay: s.delay + 's', opacity: s.opacity }"
      >✦</span>
    </template>

    <!-- Graduation caps -->
    <template v-if="mounted && kind === 'admission'">
      <span v-for="(c, i) in caps" :key="i" class="cap"
        :style="{ left: c.left + '%', animationDuration: c.duration + 's', animationDelay: c.delay + 's' }"
      >🎓</span>
    </template>

    <!-- Dust for memorial -->
    <template v-if="mounted && kind === 'memorial'">
      <span v-for="(d, i) in dust" :key="i" class="dust"
        :style="{ left: d.left + '%', top: d.top + '%', animationDuration: d.duration + 's', opacity: d.opacity }"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ themeKind?: string }>()

const kind = computed(() => (props.themeKind || 'birthday'))

const mounted = ref(false)

type Flake = { left: number; size: number; duration: number; delay: number; opacity: number }
type Drift = { left: number; top: number; duration: number; delay: number; opacity: number }
type Floaty = { left: number; duration: number; delay: number }
type Spark = { left: number; top: number; delay: number }

const flakes = ref<Flake[]>([])
const balloons = ref<Floaty[]>([])
const hearts = ref<Array<{ left: number; duration: number; delay: number; opacity: number }>>([])
const sparks = ref<Spark[]>([])
const stars = ref<Drift[]>([])
const caps = ref<Floaty[]>([])
const dust = ref<Drift[]>([])

function rnd(min: number, max: number) {
  return min + Math.random() * (max - min)
}

onMounted(() => {
  mounted.value = true

  flakes.value = Array.from({ length: 38 }).map(() => ({
    left: rnd(0, 100),
    size: rnd(2, 6),
    duration: rnd(6, 12),
    delay: rnd(-12, 0),
    opacity: rnd(0.35, 0.95),
  }))

  balloons.value = Array.from({ length: 9 }).map(() => ({
    left: rnd(0, 100),
    duration: rnd(10, 18),
    delay: rnd(-8, 0),
  }))

  hearts.value = Array.from({ length: 12 }).map(() => ({
    left: rnd(0, 100),
    duration: rnd(8, 14),
    delay: rnd(-10, 0),
    opacity: rnd(0.25, 0.75),
  }))

  sparks.value = Array.from({ length: 18 }).map(() => ({
    left: rnd(5, 95),
    top: rnd(10, 70),
    delay: rnd(0, 3),
  }))

  stars.value = Array.from({ length: 16 }).map(() => ({
    left: rnd(0, 100),
    top: rnd(0, 55),
    duration: rnd(6, 14),
    delay: rnd(0, 6),
    opacity: rnd(0.25, 0.9),
  }))

  caps.value = Array.from({ length: 10 }).map(() => ({
    left: rnd(0, 100),
    duration: rnd(9, 16),
    delay: rnd(-10, 0),
  }))

  dust.value = Array.from({ length: 24 }).map(() => ({
    left: rnd(0, 100),
    top: rnd(10, 90),
    duration: rnd(10, 22),
    delay: 0,
    opacity: rnd(0.05, 0.25),
  }))
})
</script>

<style scoped>
.vignette{
  background: radial-gradient(1200px 700px at 50% 30%, rgba(255,255,255,0.06), rgba(0,0,0,0.55));
}
.glow{
  background: radial-gradient(800px 160px at 50% 50%, rgba(255,255,255,0.14), rgba(255,255,255,0));
  filter: blur(18px);
}
.garland{
  background:
    radial-gradient(circle at 10% 60%, rgba(255,210,120,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 20% 30%, rgba(120,255,220,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 30% 60%, rgba(255,120,210,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 40% 30%, rgba(120,170,255,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 50% 60%, rgba(255,210,120,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 60% 30%, rgba(120,255,220,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 70% 60%, rgba(255,120,210,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 80% 30%, rgba(120,170,255,0.75) 0 3px, transparent 4px),
    radial-gradient(circle at 90% 60%, rgba(255,210,120,0.75) 0 3px, transparent 4px);
  opacity: .85;
  filter: drop-shadow(0 6px 12px rgba(0,0,0,.35));
  mask-image: linear-gradient(to bottom, black, black 60%, transparent);
}

.snowflake{
  position: absolute;
  top: -10px;
  border-radius: 9999px;
  background: rgba(255,255,255,0.9);
  filter: blur(0.2px);
  animation-name: snow;
  animation-timing-function: linear;
  animation-iteration-count: infinite;
}
@keyframes snow {
  0% { transform: translateY(-10px) translateX(0); }
  100% { transform: translateY(110vh) translateX(24px); }
}

.sprinkle{
  background-image:
    radial-gradient(circle, rgba(255,255,255,.16) 0 1px, transparent 2px),
    radial-gradient(circle, rgba(255,255,255,.12) 0 1px, transparent 2px);
  background-size: 18px 18px, 28px 28px;
  opacity: .55;
}

.balloon{
  position: absolute;
  bottom: -30px;
  width: 20px;
  height: 28px;
  border-radius: 9999px;
  background: rgba(255,255,255,.10);
  border: 1px solid rgba(255,255,255,.10);
  filter: blur(.1px);
  animation: floatUp linear infinite;
}
@keyframes floatUp{
  from { transform: translateY(0); opacity: .0; }
  10% { opacity: .7; }
  to { transform: translateY(-120vh); opacity: 0; }
}

.heart{
  position: absolute;
  bottom: -30px;
  font-size: 18px;
  color: rgba(255,255,255,.55);
  animation: floatUp linear infinite;
}

.spark{
  position: absolute;
  width: 2px;
  height: 18px;
  border-radius: 9999px;
  background: rgba(255,255,255,.25);
  animation: spark 1.6s ease-in-out infinite;
}
@keyframes spark {
  0%, 100% { transform: scaleY(.4); opacity: .25; }
  50% { transform: scaleY(1.2); opacity: .75; }
}

.star{
  position: absolute;
  color: rgba(255,255,255,.45);
  font-size: 18px;
  animation: twinkle 3.2s ease-in-out infinite;
}
@keyframes twinkle {
  0%,100% { transform: scale(.95); opacity: .25; }
  50% { transform: scale(1.1); opacity: .9; }
}

.cap{
  position: absolute;
  bottom: -30px;
  font-size: 18px;
  opacity: .7;
  animation: floatUp linear infinite;
}

.dust{
  position: absolute;
  width: 3px;
  height: 3px;
  border-radius: 9999px;
  background: rgba(255,255,255,.6);
  animation: drift linear infinite;
}
@keyframes drift{
  0%{ transform: translateX(0); }
  100%{ transform: translateX(60px); }
}
</style>
