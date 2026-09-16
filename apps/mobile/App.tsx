import { StatusBar } from 'expo-status-bar'
import { StyleSheet, Text, View } from 'react-native'
import { formatDateTime, formatScore } from '@devtest/shared'

export default function App() {
  const now = formatDateTime(new Date().toISOString())
  const score = formatScore(85.5)

  return (
    <View style={styles.container}>
      <Text style={styles.title}>DevTest Mobile</Text>
      <Text style={styles.info}>Now: {now}</Text>
      <Text style={styles.info}>Score: {score}</Text>
      <StatusBar style="auto" />
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  info: {
    fontSize: 16,
    marginBottom: 8,
  },
})
